package autoTrading2

import (
	"log"
	"math"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/candle"
	"upbit-api/internal/models"
)

func Reset() {

	//// 매일 9시 티커 설정
	now := time.Now()
	//
	resetTime := now.Truncate(24 * time.Hour).Add(time.Hour * 24).Add(time.Second)
	//
	setTicker := time.NewTicker(resetTime.Sub(now))

	//setTicker := time.NewTicker(time.Second)
	missingMarket := make(chan string, 200)

	for {
		select {

		// request 요청 fail over ( 업비트에서 API 요청을 모두 못받는 경우 )
		case market := <-missingMarket:
			time.Sleep(50 * time.Millisecond)
			can := candle.Market(market)

			responseDay := can.Day(count)

			if responseDay == nil {
				missingMarket <- market
				continue
			}

			cal(responseDay, market)

		case <-setTicker.C:

			log.Println("데이터 초기화 시작")

			// 데이터 초기화 실행
			// 상태값 담을 map
			PreviousMarketInfo = make(map[string]Info)

			// 지난 n일의 저점과 고점의 평균을 구함 ( 하루동안 상태값으로 남겨놓음 )

			for _, market := range config.Markets {

				//API 갯수 제한
				time.Sleep(50 * time.Millisecond)

				can := candle.Market(market)

				responseDay := can.Day(count)

				// fail over
				if responseDay == nil {
					missingMarket <- market
					continue
				}

				cal(responseDay, market)
			}

			//// 매일 9시 티커 설정
			now := time.Now()

			resetTime := now.Truncate(24 * time.Hour).Add(time.Hour * 24).Add(time.Second)

			setTicker = time.NewTicker(resetTime.Sub(now))

			log.Println("데이터 초기화 끝")

		}
	}

}

// 저점,종가,고가,퍼센트,현재가 등등 데이터 연산
func cal(responseDay *models.ResponseDay, market string) {
	lowAverage := 0.0
	highAverage := 0.0
	tradeAverage := 0.0
	openingPrice := 0.0
	tradePrice := 0.0
	for i, day := range *responseDay {
		if i == 0 {
			openingPrice += day.OpeningPrice
			tradePrice += day.TradePrice
		}

		if i != 0 {
			lowAverage += day.LowPrice
			tradeAverage += day.TradePrice
			highAverage += day.HighPrice
		}
	}
	lowAverage = lowAverage / (count - 1)
	tradeAverage = tradeAverage / (count - 1)
	highAverage = highAverage / (count - 1)
	lowTradeGap := math.Trunc((lowAverage/tradeAverage-1)*10000) / 100
	closeTradingGap := math.Trunc((tradePrice/tradeAverage-1)*10000) / 100
	highTradeGap := math.Trunc((highAverage/tradeAverage-1)*10000) / 100

	PreviousMarketMutex.Lock()
	PreviousMarketInfo[market] = Info{
		LowAverage:      lowAverage,
		TradeAverage:    tradeAverage,
		HighAverage:     highAverage,
		LowTradeGap:     lowTradeGap,
		CloseTradingGap: closeTradingGap,
		HighTradeGap:    highTradeGap,
		OpeningPrice:    openingPrice,
		TradePrice:      tradePrice,
	}
	PreviousMarketMutex.Unlock()
}
