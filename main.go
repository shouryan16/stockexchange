package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"net/http"

	"os"
	"os/signal"
	"stockexchange/pkg/engine"
	"stockexchange/pkg/handler"
	"stockexchange/pkg/server"

	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

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

	book := make(map[string]*engine.OrderBook)

	h := handler.NewStockHandler(1 * time.Second)
	s := server.NewServer(book)

	errorgroup, errorcontext := errgroup.WithContext(ctx)
	errorgroup.Go(func() error {
		return h.InitAndStart(ctx, book)
	})
	time.Sleep(time.Second * 10)
	errorgroup.Go(func() error {
		return s.Start(ctx, ":8090")
	})
	errorgroup.Go(func() error {
		return handleSignals(errorcontext, cancel)
	})
	errorgroup.Wait()
	return nil
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
