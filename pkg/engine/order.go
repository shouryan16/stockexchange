package engine

type Order struct {
	UserEmail string  `json:"UserEmail"`
	Name      string  `json:"Name"`
	Amount    int     `json:"Amount"`
	Price     float32 `json:"Price"`
	ID        string  `json:"ID"`
	Intent    string  `json:"Intent"`
	Timestamp string  `json:"Timestamp"`
}
