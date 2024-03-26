package models

import (
	"gorm.io/gorm"
	"time"
)

// 매수 체결 대기 목록 ( 9시 1분 생성후 남아 있다면 다음날 8시 55분에 모두 제거 )
type BidWaiting struct {
	Uuid   string `json:"uuid" gorm:"primaryKey;autoIncrement"`
	Ticker string `json:"ticker"`

	BidPrice  string `json:"bid_price"`
	BidVolume string `json:"bid_volume"`

	LowTradeGap     float64 `json:"low_trade_gap"`     // 저가 - 종가 차이
	CloseTradingGap float64 `json:"close_trading_gap"` // 종가평균 대비 현재가
	HighTradeGap    float64 `json:"high_trade_gap"`    // 고가 - 종가 차이

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
