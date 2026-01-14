package autoTrading2

import (
	"encoding/json"
	"log"
	"math"
	"strconv"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/connect"
	"upbit-api/internal/datastore"
	"upbit-api/internal/models"
)

// LimitOrder WebSocket으로 실시간 시세를 수신하며 조건 충족 시 지정가 매수
func LimitOrder() {
	for {
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
				tradePriceUpdate(ticker, info)

				// 시작가가 고점 평균보다 높으면 스킵 (급등 코인 제외)
				if info.HighAverage < info.OpeningPrice {
					delete(PreviousMarketInfo, ticker.Code)
					PreviousMarketMutex.Unlock()
					continue
				}

				// 시작가와 종가 평균 편차가 0.5% 이내인 경우만 진입
				if math.Abs(info.CloseTradingGap) < 0.5 {
					// 고가-저가 변동폭이 5% 미만이면 스킵 (변동성 부족)
					if (info.HighTradeGap - info.LowTradeGap) < 5 {
						delete(PreviousMarketInfo, ticker.Code)
						PreviousMarketMutex.Unlock()
						continue
					}

					// 저점 평균 가격에 지정가 매수 주문
					{
						bidPrice, bidVolume := SetBidPriceAndVolume(info)
						coin := orders.Market(ticker.Code)

						traded := coin.BidMarketLimit(bidPrice, bidVolume)

						switch traded {
						case "-1":
							log.Println("매수 주문 실패:", ticker.Code)
							delete(PreviousMarketInfo, ticker.Code)
							PreviousMarketMutex.Unlock()
							continue
						case "0": // 잔고 부족
							delete(PreviousMarketInfo, ticker.Code)
							PreviousMarketMutex.Unlock()
							continue

						default: // 매수 주문 성공
							flow := models.AutoTrading2{
								WaitingUuid:     traded,
								Ticker:          ticker.Code,
								BidPrice:        bidPrice,
								BidVolume:       bidVolume,
								BidAmount:       strconv.Itoa(OrderAmount),
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
	}

}

// tradePriceUpdate 현재가 및 종가 대비 현재가 비율 업데이트
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
