package autoTrading2

import (
	"fmt"
	"strconv"
	"sync"
)

const (
	// CandleDays 캔들 데이터 조회 일수 (4일전까지 = 3일 데이터)
	CandleDays = 4
	// OrderAmount 1회 주문 금액 (원)
	OrderAmount = 50000
)

var (
	// PreviousMarketInfo 3일 저점/종가/고가 평균 및 상태값 저장
	PreviousMarketInfo  map[string]Info
	PreviousMarketMutex sync.RWMutex
)

// Info 코인별 가격 분석 정보
type Info struct {
	LowAverage   float64 // 저가 평균
	TradeAverage float64 // 종가 평균
	HighAverage  float64 // 고가 평균

	LowTradeGap     float64 // 종가 대비 저가 차이 (%)
	CloseTradingGap float64 // 종가 평균 대비 현재가 차이 (%)
	HighTradeGap    float64 // 종가 대비 고가 차이 (%)

	OpeningPrice float64 // 금일 시작가
	TradePrice   float64 // 현재가
}

func (i Info) String() string {
	return fmt.Sprintf(
		`저점평균: %0.3f
종가평균: %0.3f
고가평균: %0.3f
금일 시작가: %f
종가평균 대비 저점: %0.2f%%
종가평균 대비 현가: %0.2f%%
종가평균 대비 고점: %0.2f%%
현재가: %f`,
		i.LowAverage, i.TradeAverage, i.HighAverage, i.OpeningPrice,
		i.LowTradeGap, i.CloseTradingGap, i.HighTradeGap, i.TradePrice)
}

// FormatPrice 업비트 호가 단위에 맞게 가격 포맷팅
// 업비트는 가격대별로 호가 단위가 다름
func FormatPrice(price float64) string {
	switch {
	case price < 1:
		return fmt.Sprintf("%0.4f", price)
	case price < 10:
		return fmt.Sprintf("%0.3f", price)
	case price < 100:
		return fmt.Sprintf("%0.2f", price)
	case price < 1000:
		return fmt.Sprintf("%0.1f", price)
	case price < 10000:
		return fmt.Sprintf("%d", int(price))
	case price < 100000:
		return fmt.Sprintf("%d", int(price)/10*10)
	case price < 1000000:
		return fmt.Sprintf("%d", int(price)/100*100)
	default:
		return fmt.Sprintf("%d", int(price)/1000*1000)
	}
}

// SetBidPriceAndVolume 매수 가격과 수량 계산
func SetBidPriceAndVolume(info Info) (string, string) {
	bidPrice := FormatPrice(info.LowAverage)

	f, err := strconv.ParseFloat(bidPrice, 64)
	if err != nil {
		panic(err)
	}
	volume := fmt.Sprintf("%0.8f", OrderAmount/f)

	return bidPrice, volume
}

// SetAskPrice 매도 가격 계산
func SetAskPrice(askFloat float64) string {
	return FormatPrice(askFloat)
}
