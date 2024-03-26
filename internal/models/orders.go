package models

import (
	"fmt"
	"time"
)

type GetOrder struct {
	Uuid string `json:"uuid"`
}

type AskOrder struct {
	Market  string `json:"market" binding:"required" example:"KRW-BTC"` // 마켓 ID
	OrdType string `json:"ord_type"`
	Side    string `json:"side" binding:"required"` // 주문 종류
	Volume  string `json:"volume"`
}

func (A AskOrder) String() string {

	return fmt.Sprintf("%s 매도 수량 : %s", A.Market, A.Volume)
}

type BidOrder struct {
	Market  string `json:"market" binding:"required" example:"KRW-BTC"` // 마켓 ID
	OrdType string `json:"ord_type"`
	Price   string `json:"price"`
	Side    string `json:"side" binding:"required"` // 주문 종류
}

func (B BidOrder) String() string {
	return fmt.Sprintf("%s 매수 : %s원", B.Market, B.Price)
}

type LimitOrder struct {
	Market  string `json:"market" binding:"required" example:"KRW-BTC"` // 마켓 ID
	OrdType string `json:"ord_type"`
	Price   string `json:"price"`
	Side    string `json:"side" binding:"required"` // 주문 종류
	Volume  string `json:"volume"`
}

func (L LimitOrder) String() string {

	var ordType string

	switch L.OrdType {
	case "bid":
		ordType = "매수"
	case "ask":
		ordType = "매도"
	}

	return fmt.Sprintf("%s: 지정가 %s 주문 주문량:%s , 주문가격:%s ", L.Market, ordType, L.Volume, L.Price)
}

type Order struct {
	Market     string `json:"market" binding:"required" example:"KRW-BTC"` // 마켓 ID
	Side       string `json:"side" binding:"required"`                     // 주문 종류
	Volume     string `json:"volume"`                                      // 주문량 (지정가, 시장가 매도 시 필수)
	Price      string `json:"price"`
	OrdType    string `json:"ord_type"`
	Identifier string `json:"identifier"`
}

type OrderList struct {
	Market      string   `json:"market"`
	Uuids       []string `json:"uuids"`
	Identifiers []string `json:"identifiers"`
	State       string   `json:"state"`
	States      []string `json:"states"`
}

type RespOrder struct {
	Uuid      string    `json:"uuid"` // 주문 고유 아이디
	Side      string    `json:"side"` // 주문 종류
	OrdType   string    `json:"ord_type"`
	Price     string    `json:"price"` // 주문 당시 화페 가격
	State     string    `json:"state"`
	Market    string    `json:"market"`
	CreatedAt time.Time `json:"created_at"`
}

type CancelOrder struct {
	//Identifier string `json:"identifier"`
	Uuid string `json:"uuid"`
}

type ResponseOrder400 struct {
	Error struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	} `json:"error"`
}
