package services

type TradeRequest struct {
	Exchange string  `json:"exchange"`
	Side     string  `json:"side"`
	Symbol   string  `json:"symbol"`
	Amount   float64 `json:"amount"`
}

type OkxInstrument struct {
	InstId string `json:"instId"`
	MinSz  string `json:"minSz"`
	LotSz  string `json:"lotSz"`
}
