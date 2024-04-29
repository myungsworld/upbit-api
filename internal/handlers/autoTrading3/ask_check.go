package autoTrading3

import (
	"log"
	"strconv"
	"time"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/datastore"
	"upbit-api/internal/middlewares"
	"upbit-api/internal/models"
)

// 매 57분마다 매도된 데이터 데이터베이스 기입
// 8시 55분에 모든 데이터를 초기화 하고 난후 8시 57분에 저장 하고 다음날 프로그램 시작하기 위해서
func AskCheck() {
	setTicker := middlewares.SetTimerEveryHourByMinute(57)
	for {
		select {
		case <-setTicker.C:

			log.Println("매도된 데이터 데이터베이스 저장 시작")
			orderList := *orders.GetDoneList()

			for _, order := range orderList {

				// 오늘인 경우
				now := time.Now().UTC()
				startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
				endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

				if order.CreatedAt.UTC().UnixNano() > startOfDay.UnixNano() && order.CreatedAt.UTC().UnixNano() < endOfDay.UnixNano() {

					// 지정가 매도가 되었다면
					if order.Side == "ask" && order.OrdType == "limit" {

						flow := models.AutoTrading2{}
						if err := datastore.DB.Model(&models.AutoTrading2{}).
							Where("ticker = ?", order.Market).
							Where("aw_created_at BETWEEN ? AND ?", startOfDay, endOfDay).
							Find(&flow).Error; err != nil {
							panic(err)
						}

						// 기입이 안되어 있을때 (중복 방지)
						if flow.AskUuid == "" && flow.Id != 0 {
							price, _ := strconv.ParseFloat(order.Price, 64)
							volume, _ := strconv.ParseFloat(order.ExecutedVolume, 64)
							fee, _ := strconv.ParseFloat(order.PaidFee, 64)

							// 매도된 금액
							askAmount := int(price*volume - fee)

							updating := map[string]interface{}{
								"ask_uuid":     order.Uuid,
								"ask_amount":   strconv.Itoa(askAmount),
								"a_created_at": order.CreatedAt,
							}

							if err := datastore.DB.Model(&flow).
								Updates(updating).Error; err != nil {
								panic(err)
							}
							log.Println(order.Market, "매도체결 데이터 데이터베이스 저장")
						}

					}

				}
			}

			log.Println("매도된 데이터 데이터베이스 저장 끝")
			setTicker = middlewares.SetTimerEveryHourByMinute(57)
		}
	}
}
