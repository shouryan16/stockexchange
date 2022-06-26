package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"stockexchange/pkg/engine"

	"github.com/rs/xid"
)

type Server struct {
	book map[string]*engine.OrderBook
}

func NewServer(book map[string]*engine.OrderBook) *Server {
	return &Server{
		book: book,
	}
}

func (s *Server) Start(ctx context.Context, port string) error {
	// http.HandleFunc("/Login", loginandsignup.Login)
	// http.HandleFunc("/Signup", loginandsignup.SignUp)
	if ctx.Err() != nil {
		log.Printf("error received while starting server: %s", ctx.Err())
		return ctx.Err()
	}
	http.HandleFunc("/order", s.Getorder)
	err := http.ListenAndServe(":"+port, nil)
	return err
}

func (s *Server) Getorder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var neworder engine.Order
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &neworder)
		if err != nil {
			log.Printf("Error parsing the body: %v", err)
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		neworder.ID = neworder.UserEmail + xid.New().String()
		if neworder.Intent == "buy" {
			s.book[neworder.Name].ProcessBuyOrder(&neworder)
		} else {
			s.book[neworder.Name].ProcessSellOrder(&neworder)
		}

		_ = json.NewEncoder(w).Encode("Order placed successfully.")

		// db := dbConn()
		// insOrder, err := db.Prepare("INSERT INTO orderhistory(ID, Order) VALUES(?,?)")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// Order, err := json.Marshal(neworder)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// insOrder.Exec(neworder.ID, Order)

		// defer db.Close()

		// context.IndentedJSON(http.StatusCreated, neworder)
	}
}
