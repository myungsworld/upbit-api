package autoTrading2

import (
	"fmt"
	"sync"
)

const count = 4

var (
	PreviousMarketInfo  map[string]Info
	PreviousMarketMutex sync.RWMutex
)

type Info struct {
	LowAverage   float64 // 저가 평균
	TradeAverage float64 // 종가 평균
	HighAverage  float64 // 고가 평균

	LowTradeGap     float64 // 저가 - 종가 차이
	TradeOpeningGap float64 // 종가 대비 시작가
	HighTradeGap    float64 // 고가 - 종가 차이

	OpeningPrice float64 // 금일 시작가

	TradePrice float64 // 현재가
}

func (i Info) String() string {

	result := fmt.Sprintf(
		`저점평균: %0.2f
종가 평균: %0.2f
고가 평균: %0.2f
금일 시작가 : %f
저점대비 종가: %0.2f%%
종가대비 시가: %0.2f%%
고점대비 종가: %0.2f%%
현재가: %f`, i.LowAverage, i.TradeAverage, i.HighAverage, i.OpeningPrice, i.LowTradeGap, i.TradeOpeningGap, i.HighTradeGap, i.TradePrice)

	return result

}
