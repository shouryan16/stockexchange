package engine

import "log"

var ledger = []*Trade{}

func RecordTransaction(trade *Trade) {
	log.Printf("%+v", trade)
	ledger = append(ledger, trade)
}
