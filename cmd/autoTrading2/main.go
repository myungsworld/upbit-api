package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/candle"
)

type info struct {
	LowAverage   float64 // 저가 평균
	TradeAverage float64 // 종가 평균
	HighAverage  float64 // 고가 평균

	LowTradeGap  float64 // 저가 - 종가 차이
	HighTradeGap float64 // 고가 - 종가 차이

	OpeningPrice float64 // 금일 시작가
}

// 이전 데이터 가져오는 기준일
const count = 4

var previousMarketInfo map[string]info

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 1일마다 리셋 ( 한국시각 9시 )
	go reset()

	oneSecTicker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-oneSecTicker.C:
				fmt.Println(len(previousMarketInfo))
			}
		}
	}()

	// opening Price 가 저점의 평균보다 낮으면 매수 안함 ( 오늘 갑자기 존나 내려간 경우 혹은 어제 많이 내려서 마감한 경우)

	// opening Price 가 고점의 평균보다 높으면 매수 안함 ( 오늘 갑자기 존나 올라온 경우 혹은 어제 많이 올라서 마감한 경우)

	// 하루전에 급격히 오른경우 다시 봐야함

	// 저점의 평균으로 왔을때 매수

	// 고점과 저점의 퍼센트 차이 구하고 그 반 퍼센트의 수익이 난다면 매도

	<-stopChan

}

func reset() {

	setTicker := time.NewTicker(time.Second)

	//
	for {
		select {
		case <-setTicker.C:

			log.Println("데이터 초기화 시작")

			// 데이터 초기화 실행
			// 상태값 담을 map
			previousMarketInfo = make(map[string]info)

			// 지난 n일의 저점과 고점의 평균을 구함 ( 하루동안 상태값으로 남겨놓음 )
			{
				// 1초에 20번

				for _, market := range config.Markets {

					// API 갯수 제한
					time.Sleep(80 * time.Millisecond)

					//if market != "KRW-GLM" {
					//	continue
					//}

					can := candle.Market(market)

					responseDay := can.Day(count)

					lowAverage := 0.0
					highAverage := 0.0
					tradeAverage := 0.0
					openingPrice := 0.0
					for i, day := range *responseDay {
						if i == 0 {
							openingPrice += day.OpeningPrice
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
					highTradeGap := math.Trunc((highAverage/tradeAverage-1)*10000) / 100

					previousMarketInfo[market] = info{
						lowAverage,
						tradeAverage,
						highAverage,
						lowTradeGap,
						highTradeGap,
						openingPrice,
					}

					fmt.Print("저점 평균:", lowAverage)
					fmt.Print("종가 평균:", tradeAverage)
					fmt.Print("고가 평균:", highAverage)
					fmt.Print("금일 시작가:", openingPrice)
					fmt.Print("저점대비 종가:", lowTradeGap)
					fmt.Println("고점대비 종가:", highTradeGap)
				}

			}
			// 매일 9시 티커 설정
			now := time.Now()

			resetTime := now.Truncate(24 * time.Hour).Add(time.Hour * 24).Add(time.Second)

			setTicker = time.NewTicker(resetTime.Sub(now))

			log.Println("데이터 초기화 끝")

		}
	}

}

func init() {
	config.Init()
}
