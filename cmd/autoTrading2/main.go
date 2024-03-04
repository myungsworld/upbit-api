package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/candle"
	"upbit-api/internal/connect"
	"upbit-api/internal/models"
)

type info struct {
	LowAverage   float64 // 저가 평균
	TradeAverage float64 // 종가 평균
	HighAverage  float64 // 고가 평균

	LowTradeGap     float64 // 저가 - 종가 차이
	TradeOpeningGap float64 // 종가 대비 시작가
	HighTradeGap    float64 // 고가 - 종가 차이

	OpeningPrice float64 // 금일 시작가
}

func (i info) String() string {

	result := fmt.Sprintf(
		`저점평균: %0.2f
종가 평균: %0.2f
고가 평균: %0.2f
금일 시작가 : %f
저점대비 종가: %0.2f%%
종가대비 시가: %0.2f%%
고점대비 종가: %0.2f%%`, i.LowAverage, i.TradeAverage, i.HighAverage, i.OpeningPrice, i.LowTradeGap, i.TradeOpeningGap, i.HighTradeGap)

	return result

}

// 이전 데이터 가져오는 기준일
const count = 4

var (
	previousMarketInfo  map[string]info
	previousMarketMutex sync.RWMutex
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 1일마다 리셋 ( 한국시각 9시 )
	go reset()

	socketOpenTicker := time.NewTicker(time.Second * 2)

	go func() {
		for {
			select {
			// 소켓 시세 수신
			case <-socketOpenTicker.C:

				conn := connect.Socket(config.Ticker)
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						break
					}

					ticker := models.Ticker{}
					if err := json.Unmarshal(message, &ticker); err != nil {
						panic(err)
					}

					previousMarketMutex.Lock()
					if data, ok := previousMarketInfo[ticker.Code]; ok {

						fmt.Println(ticker.Code)
						fmt.Println(data)

						fmt.Println(ticker.Code, "삭제")
						delete(previousMarketInfo, ticker.Code)

						// opening Price 가 저점의 평균보다 낮으면 매수 안함 ( 오늘 갑자기 존나 내려간 경우 혹은 어제 많이 내려서 마감한 경우)

						// opening Price 가 고점의 평균보다 높으면 매수 안함 ( 오늘 갑자기 존나 올라온 경우 혹은 어제 많이 올라서 마감한 경우)

						// 하루전에 급격히 오른경우 다시 봐야함

						// 저점의 평균으로 왔을때 매수

						// 고점과 저점의 퍼센트 차이 구하고 그 반 퍼센트의 수익이 난다면 매도

					}
					previousMarketMutex.Unlock()

				}

				fmt.Println("여기는 for 문 밖")

			}
		}
	}()

	//go func() {
	//	for {
	//		select {
	//		case val, ok := <-golChan:
	//			if !ok {
	//				fmt.Println(val, ok)
	//			} else {
	//				fmt.Println(val, ok)
	//				time.Sleep(time.Hour)
	//			}
	//		}
	//	}
	//}()

	// opening Price 가 저점의 평균보다 낮으면 매수 안함 ( 오늘 갑자기 존나 내려간 경우 혹은 어제 많이 내려서 마감한 경우)

	// opening Price 가 고점의 평균보다 높으면 매수 안함 ( 오늘 갑자기 존나 올라온 경우 혹은 어제 많이 올라서 마감한 경우)

	// 하루전에 급격히 오른경우 다시 봐야함

	// 저점의 평균으로 왔을때 매수

	// 고점과 저점의 퍼센트 차이 구하고 그 반 퍼센트의 수익이 난다면 매도

	<-stopChan

}

func reset() {

	//// 매일 9시 티커 설정
	//now := time.Now()
	//
	//resetTime := now.Truncate(24 * time.Hour).Add(time.Hour * 24).Add(time.Second)
	//
	//setTicker = time.NewTicker(resetTime.Sub(now))

	setTicker := time.NewTicker(time.Second)

	//
	for {
		select {
		case <-setTicker.C:

			log.Println("데이터 초기화 시작")

			// 데이터 초기화 실행
			// 상태값 담을 map
			previousMarketInfo = make(map[string]info)

			st := time.Now()

			// 지난 n일의 저점과 고점의 평균을 구함 ( 하루동안 상태값으로 남겨놓음 )
			{

				for _, market := range config.Markets {

					//if market != "KRW-SOL" && market != "KRW-GRT" {
					//	continue
					//}

					//API 갯수 제한
					time.Sleep(100 * time.Millisecond)

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
					tradeOpeningGap := math.Trunc((tradeAverage/openingPrice-1)*10000) / 100
					highTradeGap := math.Trunc((highAverage/tradeAverage-1)*10000) / 100

					previousMarketMutex.Lock()
					previousMarketInfo[market] = info{
						lowAverage,
						tradeAverage,
						highAverage,
						lowTradeGap,
						tradeOpeningGap,
						highTradeGap,
						openingPrice,
					}
					previousMarketMutex.Unlock()
				}

			}

			et := time.Now()

			fmt.Println("실행시간 :", et.Sub(st))

			//// 매일 9시 티커 설정
			//now := time.Now()
			//
			//resetTime := now.Truncate(24 * time.Hour).Add(time.Hour * 24).Add(time.Second)
			//
			//setTicker = time.NewTicker(resetTime.Sub(now))

			setTicker = time.NewTicker(time.Second * 60)
			log.Println("데이터 초기화 끝")

		}
	}

}

func init() {
	config.Init()
}
