package models

// OrderBook https://docs.upbit.com/v1.4.0/reference/websocket-orderbook
// 매도 호가 : 팔고자 하는 트레이더들이 특정 가격에 대기하는것 ( 체결량이 아님 )
type OrderBook struct {
	Ty             string          `json:"type"`
	Code           string          `json:"code"`
	TotalAskSize   float64         `json:"total_ask_size"` // 호가 매도 총 잔량
	TotalBidSize   float64         `json:"total_bid_size"` // 호가 매수 총 잔량
	OrderbookUnits []OrderBookUnit `json:"orderbook_units"`
	Timestamp      int             `json:"timestamp"`
}

type OrderBookUnit struct {
	AskPrice float64 `json:"ask_price"` // 매도 호가
	BidPrice float64 `json:"bid_price"` // 매수 호가
	AskSize  float64 `json:"ask_size"`
	BidSize  float64 `json:"bid_size"`
}
