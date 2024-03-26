package autoTrading2

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
			bidWaitings := make([]models.BidWaiting, 0)
			if err := datastore.DB.Model(&models.BidWaiting{}).
				Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).
				Find(&bidWaitings).Error; err != nil {
				panic(err)
			}

			// 매수 체결 대기 제거
			for _, bidWaiting := range bidWaitings {
				orders.Cancel(bidWaiting.Uuid)
				if err := datastore.DB.Delete(&bidWaiting).Error; err != nil {
					panic(err)
				}

			}

			log.Println("매수체결 대기 초기화 끝")

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
