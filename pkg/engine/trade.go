package engine

type Trade struct {
	Name      string  `json:"name"`
	BuyerID   string  `json:"Buyer_id"`
	SellerID  string  `json:"Seller_id"`
	Amount    int     `json:"amount"`
	Price     float32 `json:"price"`
	Timestamp string  `json:"timestamp"`
}
