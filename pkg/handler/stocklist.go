package handler

import (
	"context"
	"stockexchange/pkg/engine"
	"time"

	"golang.org/x/sync/errgroup"
)

var Stocks = []string{
	"Swiss Life AG", "Spotify", "SolarCity", "UBS AG", "SHELL", "Card Services AG", "Apple", "Samsung", "Holcim AG", "TESLA",
}

type StockHandler struct {
	stockticker time.Duration
}

func NewStockHandler(time time.Duration) *StockHandler {
	return &StockHandler{
		stockticker: time,
	}
}

func (s *StockHandler) Start(ctx context.Context, book map[string]*engine.OrderBook) error {
	errorgroup, errorcontext := errgroup.WithContext(ctx)
	for _, stock := range Stocks {
		stk := stock
		book[stock] = engine.NewOrderBook()
		errorgroup.Go(func() error {
			return s.runstockworker(stk, errorcontext)
		})
	}
	return nil
}
func (a *StockHandler) runstockworker(stock string, ctx context.Context) error {
	for {
		continue
	}

}
