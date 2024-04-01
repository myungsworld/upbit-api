package autoTrading2

import (
	"log"
	"strconv"
	"time"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/datastore"
	"upbit-api/internal/models"
)

func DeleteWaitMarket() {

	//매일 8시 55분 초기화
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 55, 00, 0, now.Location())
	duration := startTime.Sub(time.Now())
	if duration <= 0 {
		startTime = startTime.Add(24 * time.Hour)
		duration = startTime.Sub(time.Now())
	}
	deleteTicker := time.NewTicker(duration)

	//deleteTicker := time.NewTicker(time.Second * 15)

	for {
		select {
		case <-deleteTicker.C:

			log.Println("매수체결 대기 초기화 시작")

			// 그날의 모든 매수 체결 대기 가져오기
			currentTime := time.Now().UTC()
			startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)
			endOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, time.UTC)
			flow := make([]models.AutoTrading2, 0)
			if err := datastore.DB.Model(&models.AutoTrading2{}).
				Where("w_created_at BETWEEN ? AND ?", startOfDay, endOfDay).
				Find(&flow).Error; err != nil {
				panic(err)
			}

			// 매수 체결 대기 제거
			for _, trading := range flow {
				orders.Cancel(trading.WaitingUuid)
				if err := datastore.DB.Model(&trading).Update("w_deleted_at", time.Now()).Error; err != nil {
					panic(err)
				}
			}

			log.Println("매수체결 대기 초기화 끝")

			log.Println("매수 uuid 데이터 초기화 시작")
			BidLimitUuidMutex.Lock()
			BidLimitUuids = make(map[string]bool)
			BidLimitUuidMutex.Unlock()
			log.Println("매수 uuid 데이터 초기화 끝")

			time.Sleep(time.Second)

			log.Println("매도대기열 초기화 시작")
			for _, trading := range flow {
				// 매도 대기열이 있고 매도가 되지 않은 경우
				if trading.AskWaitingUuid != "" && trading.AskUuid == "" {
					orders.Cancel(trading.AskWaitingUuid)
					// 55분에 매수가 되었는데 매도가에 도달하지 않은 경우 일괄 판매 후 데이터베이스 저장
					coin := orders.Market(trading.Ticker)
					uuid := coin.AskMarketPrice(trading.ExecutedVolume)

					order := orders.Get(uuid)
					if order != nil {

						price, _ := strconv.ParseFloat(order.Price, 64)
						volume, _ := strconv.ParseFloat(order.ExecutedVolume, 64)
						fee, _ := strconv.ParseFloat(order.PaidFee, 64)

						// 매도된 금액
						askAmount := int(price*volume - fee)

						updating := map[string]interface{}{
							"ask_uuid":      uuid,
							"ask_amount":    askAmount,
							"aw_deleted_at": order.CreatedAt,
						}

						if err := datastore.DB.Model(&trading).Updates(updating).Error; err != nil {
							panic(err)
						}

					} else {
						log.Println("8시 55분 매수가 된 코인들 시장가 판매시 문제")
					}
				}
			}
			log.Println("매도대기열 초기화 끝")
			now = time.Now()
			startTime = time.Date(now.Year(), now.Month(), now.Day(), 8, 55, 00, 0, now.Location())
			duration = startTime.Sub(time.Now())
			if duration <= 0 {
				startTime = startTime.Add(24 * time.Hour)
				duration = startTime.Sub(time.Now())
			}

			deleteTicker = time.NewTicker(duration)

		}
	}

}
