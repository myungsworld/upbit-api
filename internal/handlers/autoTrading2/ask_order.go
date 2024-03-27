package autoTrading2

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/datastore"
	"upbit-api/internal/models"
)

var (
	// 매수 uuid 데이터
	BidLimitUuids     map[string]bool
	BidLimitUuidMutex sync.RWMutex
)

func AskOrder() {

	// 매수 uuid map 초기화
	BidLimitUuids = make(map[string]bool)

	// 1분마다 매수가 된 코인들 조회 후 매도 걸기
	setTicker := time.NewTicker(time.Second * 60)

	for {
		select {
		case <-setTicker.C:

			// Side 가 bid 이고 OrdType 이 limit 이며 CreatedAt 이 오늘인 경우 지정가 매도 걸기
			orderList := *orders.GetDoneList()

			for _, order := range orderList {
				if order.Side == "bid" && order.OrdType == "limit" {

					// 오늘인 경우
					now := time.Now().UTC()
					startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
					endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

					if order.CreatedAt.UTC().UnixNano() > startOfDay.UnixNano() && order.CreatedAt.UTC().UnixNano() < endOfDay.UnixNano() {

						// bidLimitUUid 는 1분마다 조회하는데 해당 map 은 8시 55분에 일괄적으로 초기화 되기 때문에 그사이 체크해야함
						// 이게 무슨 말이냐면 매수체결 대기를 걸어놨는데 매수가 된걸 내가 확인할수 없으니 1분마다 매수가 되었는지 안되었는지 확인을 해야함
						// 그래서 매수가 되었는지를 확인하는 GetDoneList() 를 1분마다 가져와서 매수가 되었다면 bidLimitUuids 에 넣고 반복되는걸 막야함
						// 이 프로그램은 8시 55분에 모든걸 초기화 하고 다음날의 스윙을 기다리기 때문에 8시 55분에 초기화가 되었는데
						// 업비트의 GetDoneList() 내역은 00시까지 기준이므로 초기화를 시켜도 또 들어오기 떄문에 8시 50분부터 9시까지 막음

						blockStart := time.Date(now.Year(), now.Month(), now.Day(), 23, 50, 0, 0, time.UTC)
						blockEnd := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
						if now.UnixNano() > blockStart.UnixNano() && now.UnixNano() < blockEnd.UnixNano() {
							log.Println("매수 uuid를 초기화 하기전 8시50분과 9시 사이라서 break")
							break
						}

						BidLimitUuidMutex.Lock()
						if _, ok := BidLimitUuids[order.Uuid]; !ok {
							BidLimitUuids[order.Uuid] = true

							// 데이터베이스에 선 저장
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
								panic(err)
							}

							fmt.Println(order.Market, "매수 체결 확인 및 데이터베이스 업데이트")

							flow := models.AutoTrading2{}
							if err := datastore.DB.Model(&models.AutoTrading2{}).
								Where("bid_uuid", order.Uuid).
								Find(&flow).Error; err != nil {
								panic(err)
							}

							if flow.Id != 0 {
								log.Println(flow.Ticker, " 매도대기열 진입")

								// 지정가 매도 생성
								// 고점과 저점의 차이 -1 퍼 만큼 오를시 매도
								lowHighGap := flow.HighTradeGap - flow.LowTradeGap - 1
								askPercent := lowHighGap/100 + 1
								bidFloat, _ := strconv.ParseFloat(flow.BidPrice, 64)
								askFloat := bidFloat * askPercent
								askPrice := SetAskPrice(askFloat)

								coin := orders.Market(flow.Ticker)
								askWaitingUuid := coin.AskMarketLimit(askPrice, flow.ExecutedVolume)

								// 데이터베이스 업데이트
								updating = map[string]interface{}{
									"ask_waiting_uuid": askWaitingUuid,
									"ask_price":        askPrice,
									"ask_volume":       flow.ExecutedVolume,
									"aw_created_at":    time.Now(),
								}
								if err := datastore.DB.Model(&flow).
									Updates(updating).Error; err != nil {
									panic(err)
								}
								log.Println(flow.Ticker, "매도 대기 데이터베이스 저장, 대기열 이탈")
							}
						}
						BidLimitUuidMutex.Unlock()

					}

				}
			}

		}
	}

}
