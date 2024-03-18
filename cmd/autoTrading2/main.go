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
	"upbit-api/internal/api/orders"
	"upbit-api/internal/connect"
	"upbit-api/internal/handlers/autoTrading2"
	"upbit-api/internal/models"
)

// 코인 스윙 자동매매
// 지나간 n일의 저점,종가,고가의 평균을 가져와 종가의 평균과 현재가의 편차가 적을때
// 저점의 평균에서 매수대기 이후 매수가 되면 고점의 평균 -1퍼에서 매도
// 하루동안 모니터링후 매수대기만 걸린 데이터들 주문리스트에서 삭제후 초기화 및 회귀
func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 1일마다 리셋 ( 한국시각 9시 1초 )
	// 3일 저점,종가,고가 평균 연산 후 상태값 저장
	go autoTrading2.Reset()

	// 매일 8시 55분 매수체결 대기가 계속 걸려 있을시 그날의 매수체결 대기 삭제
	// TODO : 이거 테스트 해봐야함
	go autoTrading2.DeleteWaitMarket()

	// 고점의 종가 -1%에서 매도
	//go autoTrading2.Handler()

	socketOpenTicker := time.NewTicker(time.Second * 2)

	go func() {
		for {
			select {

			// 초기화된 데이터 지정가 매수 체결 대기 걸기
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

							// 호가 계산

							{

								// 저점의 평균 에서 지정가 매수 체결 대기
								bidPrice, bidVolume := autoTrading2.SetBidPriceAndVolume(info)
								coin := orders.Market(ticker.Code)

								orderId := coin.BidMarketLimit(bidPrice, bidVolume)

								// 주문 ID WaitMarket map 에 상태값 저장
								if orderId != nil {

									fmt.Println(ticker.Code)
									fmt.Println(info)

									fmt.Println(ticker.Code, "매수대기 체결", bidPrice)

									delete(autoTrading2.PreviousMarketInfo, ticker.Code)
									autoTrading2.PreviousMarketMutex.Unlock()
									continue
								}

							}

							// 손절도 정해야함

						}
						// 아직 조건이 만족하지 않은 마켓들은 계속 모니터링
						//delete(autoTrading2.PreviousMarketInfo, ticker.Code)

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
