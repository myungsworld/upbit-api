package autoTrading3

import (
	"encoding/json"
	"math"
	"strconv"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/connect"
	"upbit-api/internal/datastore"
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

				// 현재가와 종가평균 대비 현재가 퍼센트 다시 업데이트
				tradePriceUpdate(ticker, info)

				// opening Price 가 고점의 평균보다 높으면 매수 안함 ( 오늘 갑자기 많이 올라온 경우 혹은 어제 많이 올라서 마감한 경우)
				if info.HighAverage < info.OpeningPrice {
					delete(PreviousMarketInfo, ticker.Code)
					PreviousMarketMutex.Unlock()
					continue
				}

				// 현재가의 평균이 마이너스 일경우
				//if info.TradeAverage < 0 {
				//
				//}

				// 금일 시작가와 종가의 평균 편차가 0.5% 내외면 모니터링 진입
				if math.Abs(info.CloseTradingGap) < 0.5 {

					// 종가 평균과 고가 평균의 차이가 5퍼 내외인 경우는 제외
					if (info.HighTradeGap - info.LowTradeGap) < 5 {
						delete(PreviousMarketInfo, ticker.Code)
						PreviousMarketMutex.Unlock()
						continue
					}

					// 호가 계산

					{

						//fmt.Println(ticker.Code)
						//fmt.Println(info)

						// 저점의 평균 에서 지정가 매수 체결 대기
						bidPrice, bidVolume := SetBidPriceAndVolume(info)
						coin := orders.Market(ticker.Code)

						traded := coin.BidMarketLimit(bidPrice, bidVolume)

						switch traded {
						// 실패
						case "-1":
							panic("에러 발생")
						// 주문 금액 부족
						case "0":
							delete(PreviousMarketInfo, ticker.Code)
							PreviousMarketMutex.Unlock()
							continue

						// 매수 성공
						default:
							// 상태값 데이터베이스 업데이트
							flow := models.AutoTrading3{
								WaitingUuid:     traded,
								Ticker:          ticker.Code,
								BidPrice:        bidPrice,
								BidVolume:       bidVolume,
								BidAmount:       strconv.Itoa(amount),
								LowTradeGap:     info.LowTradeGap,
								CloseTradingGap: info.CloseTradingGap,
								HighTradeGap:    info.HighTradeGap,
								WCreatedAt:      time.Now(),
							}

							if err = datastore.DB.Create(&flow).Error; err != nil {
								panic(err)
							}

							delete(PreviousMarketInfo, ticker.Code)
							PreviousMarketMutex.Unlock()
							continue
						}

					}

				}

			}
			PreviousMarketMutex.Unlock()

		}

		// for 문 밖

	}

}

func tradePriceUpdate(ticker models.Ticker, info Info) {
	newCloseTradingGap := math.Trunc((ticker.TradePrice/info.TradeAverage-1)*10000) / 100
	PreviousMarketInfo[ticker.Code] = Info{
		LowAverage:      info.LowAverage,
		TradeAverage:    info.TradeAverage,
		HighAverage:     info.HighAverage,
		LowTradeGap:     info.LowTradeGap,
		CloseTradingGap: newCloseTradingGap,
		HighTradeGap:    info.HighTradeGap,
		OpeningPrice:    info.OpeningPrice,
		TradePrice:      ticker.TradePrice,
	}
}
