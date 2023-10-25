package models

type OrderBook struct {
	//Ty             string          `json:"type"`
	Code         string  `json:"code"`
	TotalAskSize float64 `json:"total_ask_size"` // 호가 매도 총 잔량
	//TotalBidSize   float64         `json:"total_bid_size"` // 호가 매수 총 잔량
	OrderbookUnits []OrderBookUnit `json:"orderbook_units"`
	//Timestamp      int             `json:"timestamp"`
}

type OrderBookUnit struct {
	//AskPrice float64 `json:"ask_price"` // 매도 호가
	//BidPrice float64 `json:"bid_price"` // 매수 호가
	AskSize float64 `json:"ask_size"`
	//BidSize float64 `json:"bid_size"`
}
