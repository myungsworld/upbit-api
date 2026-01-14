package autoTrading2

import (
	"log"
	"math"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/candle"
	"upbit-api/internal/models"
)

// Reset 매일 09:00:01 KST에 전체 코인의 3일 평균 데이터 초기화
func Reset() {
	setTicker := time.NewTicker(time.Second)
	missingMarket := make(chan string, 200)

	for {
		select {

		// API 요청 실패 시 재시도
		case market := <-missingMarket:
			time.Sleep(50 * time.Millisecond)
			can := candle.Market(market)

			responseDay := can.Day(CandleDays)

			if responseDay == nil {
				missingMarket <- market
				continue
			}

			cal(responseDay, market)

		case <-setTicker.C:

			log.Println("데이터 초기화 시작")

			// 상태값 초기화 및 3일 평균 계산
			PreviousMarketInfo = make(map[string]Info)

			for _, market := range config.Markets {
				time.Sleep(50 * time.Millisecond) // API 요청 제한 준수

				can := candle.Market(market)

				responseDay := can.Day(CandleDays)

				if responseDay == nil {
					missingMarket <- market
					continue
				}

				cal(responseDay, market)
			}

			// 다음 09:00:01 까지 대기
			now := time.Now()
			resetTime := now.Truncate(24 * time.Hour).Add(24*time.Hour + time.Second)
			setTicker = time.NewTicker(resetTime.Sub(now))

			log.Println("데이터 초기화 끝")

		}
	}

}

// cal 저점/종가/고가 평균 및 변동률 계산
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
	lowAverage = lowAverage / (CandleDays - 1)
	tradeAverage = tradeAverage / (CandleDays - 1)
	highAverage = highAverage / (CandleDays - 1)
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
