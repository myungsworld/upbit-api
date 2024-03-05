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

	TradePrice float64 // 현재가
}

func (i info) String() string {

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

// 이전 데이터 가져오는 기준일
const count = 4

var (
	previousMarketInfo  map[string]info
	previousMarketMutex sync.RWMutex
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 1일마다 리셋 ( 한국시각 9시 1초 )
	// 3일 저점,종가,고가 평균 연산 후 상태값 저장
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

						// opening Price(시작가) 가 저점의 평균보다 낮으면 매수 안함 ( 오늘 갑자기 존나 내려간 경우 혹은 어제 많이 내려서 마감한 경우)
						// opening Price 가 고점의 평균보다 높으면 매수 안함 ( 오늘 갑자기 존나 올라온 경우 혹은 어제 많이 올라서 마감한 경우)
						if data.LowAverage > data.OpeningPrice || data.HighAverage < data.OpeningPrice {
							delete(previousMarketInfo, ticker.Code)
							previousMarketMutex.Unlock()
							continue
						}

						// 금일 시작가와 종가의 평균 편차가 0.5% 내외면 모니터링 진입 (고루틴)
						if math.Abs(data.TradeOpeningGap) < 0.5 {

							// 종가 평균과 고가 평균의 차이가 3퍼 내외인 경우는 제외
							if (data.HighTradeGap - data.LowTradeGap) < 3 {
								delete(previousMarketInfo, ticker.Code)
								previousMarketMutex.Unlock()
								continue
							}

							fmt.Println(ticker.Code)
							fmt.Println(data)

							// 저점의 평균에서 매수
							// 고점의 종가 -1%에서 매도
							//go autoTrading2.Handler()

						}

						delete(previousMarketInfo, ticker.Code)

					}
					previousMarketMutex.Unlock()

				}

				// for 문 밖

			}
		}

		// goroutine 밖
	}()

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

					//if market != "KRW-MANA" && market != "KRW-CRO" {
					//	continue
					//}

					//API 갯수 제한
					time.Sleep(100 * time.Millisecond)

					can := candle.Market(market)

					responseDay := can.Day(count)

					if responseDay == nil {
						continue
					}

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
					tradeOpeningGap := math.Trunc((openingPrice/tradeAverage-1)*10000) / 100
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
						tradePrice,
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
