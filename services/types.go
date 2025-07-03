package services

type TradeRequest struct {
	Exchange string `json:"exchange"`
	Side     string `json:"side"`
	Symbol   string `json:"symbol"`
	Amount   int64  `json:"amount"`
}
