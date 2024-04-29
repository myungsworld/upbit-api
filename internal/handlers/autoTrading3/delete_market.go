package autoTrading3

import (
	"log"
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
