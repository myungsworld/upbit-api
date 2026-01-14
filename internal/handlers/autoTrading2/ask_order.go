package autoTrading2

import (
	"log"
	"strconv"
	"sync"
	"time"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/datastore"
	"upbit-api/internal/models"
)

var (
	BidLimitUuids     map[string]bool // 처리된 매수 UUID 추적
	BidLimitUuidMutex sync.RWMutex
)

// AskOrder 1분마다 매수 체결 확인 후 매도 주문 설정
func AskOrder() {
	BidLimitUuids = make(map[string]bool)
	setTicker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-setTicker.C:
			orderList := *orders.GetDoneList()

			for _, order := range orderList {
				// 오늘 체결된 지정가 매수 주문만 처리
				if order.Side != "bid" || order.OrdType != "limit" {
					continue
				}

				now := time.Now().UTC()
				startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
				endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

				if order.CreatedAt.UTC().UnixNano() <= startOfDay.UnixNano() ||
					order.CreatedAt.UTC().UnixNano() >= endOfDay.UnixNano() {
					continue
				}

				// 08:50 ~ 09:00 사이에는 처리 중단 (일일 초기화 시간대)
				blockStart := time.Date(now.Year(), now.Month(), now.Day(), 23, 50, 0, 0, time.UTC)
				blockEnd := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
				if now.After(blockStart) && now.Before(blockEnd) {
					log.Println("일일 초기화 시간대 - 매도 주문 설정 대기")
					break
				}

				BidLimitUuidMutex.Lock()
				if _, ok := BidLimitUuids[order.Uuid]; !ok {
					BidLimitUuids[order.Uuid] = true

					// 매수 체결 정보 DB 저장
					updating := map[string]interface{}{
						"bid_uuid":        order.Uuid,
						"price":           order.Price,
						"request_volume":  order.Volume,
						"executed_volume": order.ExecutedVolume,
						"b_created_at":    order.CreatedAt,
					}

					if err := datastore.DB.Model(&models.AutoTrading2{}).
						Where("ticker = ?", order.Market).
						Where("w_created_at BETWEEN ? AND ?", startOfDay, endOfDay).
						Where("b_created_at is null").
						Updates(updating).Error; err != nil {
						log.Println("매수 체결 DB 저장 실패:", err)
						BidLimitUuidMutex.Unlock()
						continue
					}

					log.Println(order.Market, "매수 체결 확인")

					flow := models.AutoTrading2{}
					if err := datastore.DB.Model(&models.AutoTrading2{}).
						Where("bid_uuid", order.Uuid).
						Find(&flow).Error; err != nil {
						log.Println("매수 데이터 조회 실패:", err)
						BidLimitUuidMutex.Unlock()
						continue
					}

					if flow.Id != 0 && flow.AskWaitingUuid == "" {
						log.Println(flow.Ticker, "매도 대기열 진입")

						// 고가-저가 차이의 -1% 수준에서 매도
						lowHighGap := flow.HighTradeGap - flow.LowTradeGap - 1
						askPercent := lowHighGap/100 + 1
						bidFloat, _ := strconv.ParseFloat(flow.BidPrice, 64)
						askFloat := bidFloat * askPercent
						askPrice := SetAskPrice(askFloat)

						coin := orders.Market(flow.Ticker)
						askWaitingUuid := coin.AskMarketLimit(askPrice, flow.ExecutedVolume)

						// 매도 대기 정보 DB 저장
						updating = map[string]interface{}{
							"ask_waiting_uuid": askWaitingUuid,
							"ask_price":        askPrice,
							"ask_volume":       flow.ExecutedVolume,
							"aw_created_at":    time.Now(),
						}
						if err := datastore.DB.Model(&flow).Updates(updating).Error; err != nil {
							log.Println("매도 대기 DB 저장 실패:", err)
						}
						log.Println(flow.Ticker, "매도 대기 설정 완료")
					}
				}
				BidLimitUuidMutex.Unlock()
			}
		}
	}
}
