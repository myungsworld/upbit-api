package models

import "time"

type AutoTrading3 struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Ticker string `json:"ticker"`

	// 매수대기 데이터
	WaitingUuid     string    `json:"waiting_uuid"`
	BidPrice        string    `json:"bid_price"`
	BidVolume       string    `json:"bid_volume"`
	BidAmount       string    `json:"bid_amount"`
	LowTradeGap     float64   `json:"low_trade_gap"`     // 저가 - 종가 차이
	CloseTradingGap float64   `json:"close_trading_gap"` // 종가평균 대비 현재가
	HighTradeGap    float64   `json:"high_trade_gap"`    // 고가 - 종가 차이
	WCreatedAt      time.Time `json:"w_created_at"`
	WDeletedAt      time.Time `json:"w_deleted_at" gorm:"default:null;"`

	// 매수 데이터
	BidUuid        string    `json:"bid_uuid"`
	Price          string    `json:"price"`           // 주문 당시 화폐 가격
	RequestVolume  string    `json:"request_volume"`  // 사용자가 입력한 주문 양
	ExecutedVolume string    `json:"executed_volume"` // 실제 체결된 양
	BCreatedAt     time.Time `json:"b_created_at" gorm:"default:null;"`
	//BDeletedAt     time.Time `json:"b_deleted_at" gorm:"default:null;"`

	// 매도대기 데이터
	AskWaitingUuid string    `json:"ask_waiting_uuid"`
	AskPrice       string    `json:"ask_price"`
	AskVolume      string    `json:"ask_volume"`
	AwCreatedAt    time.Time `json:"aw_created_at" gorm:"default:null;"`
	AwDeletedAt    time.Time `json:"aw_deleted_at" gorm:"default:null;"`

	// 매도 데이터
	AskUuid    string    `json:"ask_uuid"`
	AskAmount  string    `json:"ask_amount"`
	ACreatedAt time.Time `json:"a_created_at" gorm:"default:null;"`
}
