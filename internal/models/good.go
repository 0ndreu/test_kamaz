package models

type Good struct {
	Price    int64 `json:"price"`
	Quantity int64 `json:"quantity"`
	Amount   int64 `json:"amount"`
	Object   int64 `json:"object"`
	Method   int64 `json:"method"`
}
