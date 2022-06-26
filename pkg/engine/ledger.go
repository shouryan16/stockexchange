package engine

import "fmt"

var ledger = []*Trade{}

func RecordTransaction(trade *Trade) {
	fmt.Printf("%v", trade)
	ledger = append(ledger, trade)
}
 