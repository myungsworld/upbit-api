package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/datastore"
	"upbit-api/internal/models"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

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
						"a_created_at": time.Now(),
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

	<-stopChan

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
	datastore.ConnectDB()
}
