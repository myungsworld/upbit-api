package models

type ResponseDay []struct {
	Market            string  `json:"market"`
	CandleDateTimeKst string  `json:"candle_date_time_kst"`
	OpeningPrice      float64 `json:"opening_price"`
	HighPrice         float64 `json:"high_price"`
	LowPrice          float64 `json:"low_price"`
	TradePrice        float64 `json:"trade_price"` // 종가
}
