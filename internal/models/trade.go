package models

// Trade
// https://docs.upbit.com/v1.4.0/reference/websocket-trade
type Trade struct {
	Ty               string  `json:"ty"`                 // trade : 체결
	Code             string  `json:"code"`               // 마켓 코드 (ex.KRW-BTC)
	TradePrice       float64 `json:"trade_price"`        // 체결가격
	TradeVolume      float64 `json:"trade_volume"`       // 체결량
	AskBid           string  `json:"ask_bid"`            // 매수,매도 구분 ASK:매도, BID:매수
	PrevClosingPrice float64 `json:"prev_closing_price"` // 전일 종가
	Change           string  `json:"change"`             // 전일대비 RISE:상승,EVEN:보합,FALL:하락
	ChangePrice      float64 `json:"change_price"`       // 부호 없는 전일 대비 값
	StreamType       string  `json:"stream_type"`        // 스트림 타입 SNAPSHOT:스냅샷,REALTIME:실시간

	//TradeDate        string  `json:"trade_date"`         // 체결 일자(UTC 기준)
	//TradeTime        string  `json:"trade_time"`         // 체결시각(UTC 기준)
	//TradeTimestamp int64 `json:"trade_timestamp"` // 체결 타임 스탬프
	//Timestamp int64 `json:"timestamp"` // 타임 스탬프
	//SequentialId int64 `json:"sequential_id"` // 체결 번호

}
