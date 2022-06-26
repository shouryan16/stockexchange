package engine

import (
	"math"
	"time"

	aq "github.com/emirpasic/gods/queues/arrayqueue"
)

type OrderBook struct {
	BuyOrders  map[float32]*aq.Queue
	SellOrders map[float32]*aq.Queue
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		BuyOrders:  make(map[float32]*aq.Queue),
		SellOrders: make(map[float32]*aq.Queue),
	}
}

func (book *OrderBook) ProcessBuyOrder(order *Order) {
	for order.Amount > 0 {
		item, isEmpty := book.SellOrders[order.Price].Peek()
		if !isEmpty {
			break
		}
		sellOrder := item.(*Order)
		minAmount := math.Min(float64(sellOrder.Amount), float64(order.Amount))
		sellOrder.Amount -= int(minAmount)
		order.Amount -= int(minAmount)
		if order.Amount == 0 {
			book.SellOrders[order.Price].Dequeue()
		}
		trade := &Trade{
			Name:      order.Name,
			BuyerID:   order.ID,
			SellerID:  sellOrder.ID,
			Amount:    int(minAmount),
			Price:     order.Price,
			Timestamp: time.Now().String(),
		}

		RecordTransaction(trade)
	}
	if order.Amount > 0 {
		book.BuyOrders[order.Price].Enqueue(&order)
	}
}

func (book *OrderBook) ProcessSellOrder(order *Order) {
	for order.Amount > 0 {
		item, isEmpty := book.BuyOrders[order.Price].Peek()
		if !isEmpty {
			break
		}
		buyOrder := item.(*Order)
		minAmount := math.Min(float64(buyOrder.Amount), float64(order.Amount))
		buyOrder.Amount -= int(minAmount)
		order.Amount -= int(minAmount)
		if order.Amount == 0 {
			book.BuyOrders[order.Price].Dequeue()
		}
		trade := &Trade{
			Name:      order.Name,
			BuyerID:   buyOrder.ID,
			SellerID:  order.ID,
			Amount:    int(minAmount),
			Price:     order.Price,
			Timestamp: time.Now().String(),
		}
		RecordTransaction(trade)
	}
	if order.Amount > 0 {
		book.SellOrders[order.Price].Enqueue(&order)
	}
}
