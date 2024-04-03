package autoTrading2

import (
	"log"
	"strconv"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/datastore"
	"upbit-api/internal/middlewares"
	"upbit-api/internal/models"
)

func UpdateDB() {

	sixHourTicker := middlewares.SetTimerEvery6Hour()

	for {
		select {
		case <-sixHourTicker.C:

			getOrders := *orders.GetDoneList()

			for _, order := range getOrders {
				// 시장가 매도인 경우 중
				if order.Side == "ask" && order.OrdType == "limit" {
					flow := models.AutoTrading2{}
					if err := datastore.DB.Model(&models.AutoTrading2{}).
						Where("ticker", order.Market).
						Where("executed_volume", order.ExecutedVolume).
						Find(&flow).Error; err != nil {
						panic(err)
					}

					// 해당되는 flow 가 있고 매도대기데이터가 있는 경우 매도 업데이트
					if flow.Id != 0 && flow.AskWaitingUuid != "" && flow.AskUuid == "" {
						log.Println(order.Market, "드디어 매도 체결!! 쏴리질럿~~~")
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

					}

				}
			}
			sixHourTicker = middlewares.SetTimerEvery6Hour()
		}
	}

}
