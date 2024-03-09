package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
	"upbit-api/config"
	"upbit-api/internal/connect"
	"upbit-api/internal/handlers/autoTrading2"
	"upbit-api/internal/models"
)

// 이전 데이터 가져오는 기준일
const count = 4

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 1일마다 리셋 ( 한국시각 9시 1초 )
	// 3일 저점,종가,고가 평균 연산 후 상태값 저장
	go autoTrading2.Reset()

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

					autoTrading2.PreviousMarketMutex.Lock()
					if info, ok := autoTrading2.PreviousMarketInfo[ticker.Code]; ok {

						// opening Price(시작가) 가 저점의 평균보다 낮으면 매수 안함 ( 오늘 갑자기 많이 내려간 경우 혹은 어제 많이 내려서 마감한 경우)
						// opening Price 가 고점의 평균보다 높으면 매수 안함 ( 오늘 갑자기 많이 올라온 경우 혹은 어제 많이 올라서 마감한 경우)
						if info.LowAverage > info.OpeningPrice || info.HighAverage < info.OpeningPrice {
							delete(autoTrading2.PreviousMarketInfo, ticker.Code)
							autoTrading2.PreviousMarketMutex.Unlock()
							continue
						}

						// 금일 시작가와 종가의 평균 편차가 0.5% 내외면 모니터링 진입
						if math.Abs(info.CloseTradingGap) < 0.5 {

							// 종가 평균과 고가 평균의 차이가 3퍼 내외인 경우는 제외
							if (info.HighTradeGap - info.LowTradeGap) < 3 {
								delete(autoTrading2.PreviousMarketInfo, ticker.Code)
								autoTrading2.PreviousMarketMutex.Unlock()
								continue
							}

							// 저점의 평균에서 매수
							// 호가 계산

							fmt.Println(ticker.Code)
							fmt.Println(info)

							//if ticker.Code == "KRW-STORJ" {
							//	fmt.Println(info)
							//	fmt.Println("진입")
							//	bidPrice, bidVolume := autoTrading2.SetBidPriceAndVolume(info)
							//	coin := orders.Market(ticker.Code)
							//	coin.BidMarketLimit(bidPrice, bidVolume)
							//	fmt.Println(bidPrice, bidVolume)
							//}

							// 고점의 종가 -1%에서 매도
							//go autoTrading2.Handler()

							// 손절도 정해야함

						}

						//

						delete(autoTrading2.PreviousMarketInfo, ticker.Code)

					}
					autoTrading2.PreviousMarketMutex.Unlock()

				}

				// for 문 밖

			}
		}

		// goroutine 밖
	}()

	<-stopChan

}

func init() {
	config.Init()
}
