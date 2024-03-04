package models

import (
	"fmt"
)

// Ticker 현재가
// https://docs.upbit.com/reference/websocket-ticker
type Ticker struct {
	Ty                 string  `json:"type"`
	Code               string  `json:"code"`
	OpeningPrice       float64 `json:"opening_price"`
	HighPrice          float64 `json:"high_price"`
	LowPrice           float64 `json:"low_price"`
	TradePrice         float64 `json:"trade_price"` // 현재가
	PrevClosingPrice   float64 `json:"prev_closing_price"`
	Change             string  `json:"change"`
	ChangePrice        float64 `json:"change_price"`
	SignedChangePrice  float64 `json:"signed_change_price"`
	ChangeRate         float64 `json:"change_rate"`
	SignedChangeRate   float64 `json:"signed_change_rate"`
	TradeVolume        float64 `json:"trade_volume"`
	AccTradeVolume     float64 `json:"acc_trade_volume"`
	AccTradeVolume24h  float64 `json:"acc_trade_volume_24h"`
	AccTradePrice      float64 `json:"acc_trade_price"`
	AccTradePrice24h   float64 `json:"acc_trade_price_24h"`
	TradeDate          string  `json:"trade_date"`
	TradeTime          string  `json:"trade_time"`
	TradeTimestamp     int64   `json:"trade_timestamp"`
	AskBid             string  `json:"ask_bid"`
	AccAskVolume       float64 `json:"acc_ask_volume"`
	AccBidVolume       float64 `json:"acc_bid_volume"`
	Highest52WeekPrice float64 `json:"highest_52_week_price"`
	Highest52WeekDate  string  `json:"highest_52_week_date"`
	Lowest52WeekPrice  float64 `json:"lowest_52_week_price"`
	Lowest52WeekDate   string  `json:"lowest_52_week_date"`
	TradeStatus        string  `json:"trade_status"`
	MarketState        string  `json:"market_state"`
	MarketStateForIOS  string  `json:"market_state_for_ios"`
	IsTradingSuspended bool    `json:"is_trading_suspended"`
	//DelistingDate      time.Time `json:"delisting_date"`
	MarketWarning string `json:"market_warning"`
	Timestamp     int64  `json:"timestamp"`
	StreamType    string `json:"stream_type"`
}

func (t Ticker) String() string {

	var 억 uint16
	var 만 uint16
	var 원 uint16

	if t.TradePrice/1e8 > 1 {
		억 = uint16(t.TradePrice / 1e8)
		만 = uint16(t.TradePrice/1e4) % 10000
		원 = uint16(int(t.TradePrice) % 10000)
		return fmt.Sprintf("%v: 현재가:%d억%d만%d원 전일대비:%0.f퍼", t.Code, 억, 만, 원, t.SignedChangeRate)
	} else {
		만 = uint16(t.TradePrice/1e4) % 10000
		원 = uint16(int(t.TradePrice) % 10000)

		return fmt.Sprintf("%v: 현재가:%d만%d원 전일대비:%0.2f%%", t.Code, 만, 원, 100*t.SignedChangeRate)
	}

}
