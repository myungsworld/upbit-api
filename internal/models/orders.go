package models

type AskOrder struct {
	Market  string `json:"market" binding:"required" example:"KRW-BTC"` // 마켓 ID
	OrdType string `json:"ord_type"`
	Side    string `json:"side" binding:"required"` // 주문 종류
	Volume  string `json:"volume"`
}

type BidOrder struct {
	Market  string `json:"market" binding:"required" example:"KRW-BTC"` // 마켓 ID
	OrdType string `json:"ord_type"`
	Price   string `json:"price"`
	Side    string `json:"side" binding:"required"` // 주문 종류
}

type Order struct {
	Market     string `json:"market" binding:"required" example:"KRW-BTC"` // 마켓 ID
	Side       string `json:"side" binding:"required"`                     // 주문 종류
	Volume     string `json:"volume"`                                      // 주문량 (지정가, 시장가 매도 시 필수)
	Price      string `json:"price"`
	OrdType    string `json:"ord_type"`
	Identifier string `json:"identifier"`
}

type RespOrder struct {
	Uuid    string `json:"uuid"` // 주문 고유 아이디
	Side    string `json:"side"` // 주문 종류
	OrdType string `json:"ord_type"`
}
