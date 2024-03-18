package autoTrading2

import (
	"fmt"
	"strconv"
	"sync"
)

const (
	// 몇일 기준 ( 4일전까지의 데이터 = 3일의 데이터를 가져와서 금일과 비교 )
	count = 4
	// 주문 금액
	amount = 50000
)

var (
	// PreviousMarketInfo 3일 저점,종가,고가 평균 연산 후 상태값 저장
	PreviousMarketInfo  map[string]Info
	PreviousMarketMutex sync.RWMutex
)

type Info struct {
	LowAverage   float64 // 저가 평균
	TradeAverage float64 // 종가 평균
	HighAverage  float64 // 고가 평균

	LowTradeGap     float64 // 저가 - 종가 차이
	CloseTradingGap float64 // 종가평균 대비 현재가
	HighTradeGap    float64 // 고가 - 종가 차이

	OpeningPrice float64 // 금일 시작가

	TradePrice float64 // 현재가
}

func (i Info) String() string {

	result := fmt.Sprintf(
		`저점평균: %0.3f
종가 평균: %0.3f
고가 평균: %0.3f
금일 시작가 : %f
종가평균 대비 저점 : %0.2f%%
종가평균 대비 현가 : %0.2f%%
종가평균 대비 고점 : %0.2f%%
현재가: %f`, i.LowAverage, i.TradeAverage, i.HighAverage, i.OpeningPrice, i.LowTradeGap, i.CloseTradingGap, i.HighTradeGap, i.TradePrice)

	return result

}

// 지정가 매수 매도 호가 측정
func SetBidPriceAndVolume(info Info) (string, string) {

	tradePrice := info.TradePrice
	bidFloat := info.LowAverage

	var bidPrice string

	switch {
	case tradePrice < 1:
		// TODO : 1원 미만일떄 뭐해야함?
	// 현재가가 1원과 10원 사이인 코인
	case 1 < tradePrice && tradePrice < 10:
		bidPrice = fmt.Sprintf("%0.3f", bidFloat)

	// 현재가가 10원과 100원 사이인 코인
	case 10 <= tradePrice && tradePrice < 100:
		switch {
		// 저점 평균이 1원과 10원 사이인 경우
		case 1 < bidFloat && bidFloat < 10:
			bidPrice = fmt.Sprintf("%0.3f", bidFloat)
		default:
			bidPrice = fmt.Sprintf("%0.2f", bidFloat)
		}
	// 현재가가 100원과 1000원 사이인 코인
	case 100 <= tradePrice && tradePrice < 1000:
		switch {
		// 저점 평균이 10원과 100원 사이인 경우
		case 10 < bidFloat && bidFloat < 100:
			bidPrice = fmt.Sprintf("%0.2f", bidFloat)
		default:
			bidPrice = fmt.Sprintf("%0.1f", bidFloat)
		}

	// 현재가가 1000원과 10000원 사이인 코인
	case 1000 <= tradePrice && tradePrice < 10000:
		switch {
		// 저점 평균이 100원과 1000원 사이인 경우
		case 100 < bidFloat && bidFloat < 1000:
			bidPrice = fmt.Sprintf("%0.1f", bidFloat)
		default:
			bidPrice = fmt.Sprintf("%d", int(bidFloat))
		}
	// 현재가가 10000원과 50000원 사이인 코인
	case 10000 <= tradePrice && tradePrice < 100000:
		// 저점 평균이 1000원과 10000원사인 경우
		switch {
		case 1000 < bidFloat && bidFloat < 10000:
			bidPrice = fmt.Sprintf("%d", int(bidFloat))
		default:
			bidPrice = fmt.Sprintf("%d", int(bidFloat)/10*10)
		}

	// 현재가가 50000원과 1000000원 사이인 코인
	case 100000 <= tradePrice && tradePrice < 1000000:
		switch {
		case 10000 < bidFloat && bidFloat < 100000:
			bidPrice = fmt.Sprintf("%d", int(bidFloat)/10*10)
		default:
			bidPrice = fmt.Sprintf("%d", int(bidFloat)/100*100)
		}
	case 1000000 <= tradePrice:
		switch {
		case 100000 < bidFloat && bidFloat < 1000000:
			bidPrice = fmt.Sprintf("%d", int(bidFloat)/100*100)
		default:
			bidPrice = fmt.Sprintf("%d", int(bidFloat)/1000*1000)
		}
	}

	f, err := strconv.ParseFloat(bidPrice, 64)
	if err != nil {
		panic(err)
	}
	volume := fmt.Sprintf("%0.8f", amount/f)

	return bidPrice, volume
}
