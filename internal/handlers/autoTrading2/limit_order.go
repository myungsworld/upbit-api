package autoTrading2

import (
	"encoding/json"
	"fmt"
	"math"
	"upbit-api/config"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/connect"
	"upbit-api/internal/models"
)

func LimitOrder() {

	for {

		// 초기화된 데이터 지정가 매수 체결 대기 걸기
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

			PreviousMarketMutex.Lock()
			if info, ok := PreviousMarketInfo[ticker.Code]; ok {

				// opening Price(시작가) 가 저점의 평균보다 낮으면 매수 안함 ( 오늘 갑자기 많이 내려간 경우 혹은 어제 많이 내려서 마감한 경우)
				// opening Price 가 고점의 평균보다 높으면 매수 안함 ( 오늘 갑자기 많이 올라온 경우 혹은 어제 많이 올라서 마감한 경우)
				if info.LowAverage > info.OpeningPrice || info.HighAverage < info.OpeningPrice {
					delete(PreviousMarketInfo, ticker.Code)
					PreviousMarketMutex.Unlock()
					continue
				}

				// 금일 시작가와 종가의 평균 편차가 0.5% 내외면 모니터링 진입
				if math.Abs(info.CloseTradingGap) < 0.5 {

					// 종가 평균과 고가 평균의 차이가 3퍼 내외인 경우는 제외
					if (info.HighTradeGap - info.LowTradeGap) < 3 {
						delete(PreviousMarketInfo, ticker.Code)
						PreviousMarketMutex.Unlock()
						continue
					}

					// 호가 계산

					{

						// 저점의 평균 에서 지정가 매수 체결 대기
						bidPrice, bidVolume := SetBidPriceAndVolume(info)
						coin := orders.Market(ticker.Code)

						traded := coin.BidMarketLimit(bidPrice, bidVolume)

						if traded {

							fmt.Println(ticker.Code)
							fmt.Println(info)

							delete(PreviousMarketInfo, ticker.Code)
							PreviousMarketMutex.Unlock()
							continue
						}

					}

				}
				// 아직 조건이 만족하지 않은 마켓들은 계속 모니터링
				//delete(autoTrading2.PreviousMarketInfo, ticker.Code)

			}
			PreviousMarketMutex.Unlock()

		}

		// for 문 밖

	}

}
