package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"os"
	"os/signal"
	"stockexchange/pkg/engine"
	"stockexchange/pkg/handler"
	"stockexchange/pkg/loginandsignup"

	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/sync/errgroup"
)

var book = make(map[string]*engine.OrderBook)

func main() {

	if err := run(); err != nil {
		switch err {
		case context.Canceled:
			log.Fatal("context was canceled")
		case http.ErrServerClosed:
			log.Fatal("server close error")
		default:
			log.Fatalf("cannot run service because: %v", err)
		}

	}

}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := handler.NewStockHandler(1 * time.Second)

	errorgroup, errorcontext := errgroup.WithContext(ctx)
	errorgroup.Go(func() error {
		return h.Start(ctx, book)
	})
	time.Sleep(time.Second * 10)
	errorgroup.Go(func() error {

		return handleSignals(errorcontext, cancel)
	})
	startServer()
	errorgroup.Wait()
	return nil
}

func startServer() {
	router := gin.Default()
	router.POST("/Login", loginandsignup.Login)
	router.POST("/Signup", loginandsignup.SignUp)
	router.POST("/Order", Getorder)
	router.Run("localhost:9090")
}

func Getorder(context *gin.Context) {
	var neworder engine.Order
	if err := context.BindJSON(&neworder); err != nil {
		return
	}
	neworder.ID = neworder.UserEmail + xid.New().String()
	if neworder.Intent == "buy" {
		book[neworder.Name].ProcessBuyOrder(&neworder)
	} else {
		book[neworder.Name].ProcessSellOrder(&neworder)
	}

	db := dbConn()
	insOrder, err := db.Prepare("INSERT INTO orderhistory(ID, Order) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	Order, err := json.Marshal(neworder)
	if err != nil {
		fmt.Println(err)
	}

	insOrder.Exec(neworder.ID, Order)

	defer db.Close()

	context.IndentedJSON(http.StatusCreated, neworder)

}

func dbConn() (db *sql.DB) {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/userdetails")
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return db
}

func handleSignals(ctx context.Context, cancel context.CancelFunc) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-sigCh:
		log.Printf("got signal %v, stopping", s)
		cancel()
		return nil
	case <-ctx.Done():
		log.Printf("context is done")
		return ctx.Err()
	}
}
