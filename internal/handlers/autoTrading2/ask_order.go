package autoTrading2

import (
	"time"
	"upbit-api/internal/datastore"
	"upbit-api/internal/models"
)

func AskOrder() {

	setTicker := time.NewTicker(time.Second)

	for {
		select {
		case <-setTicker.C:

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

			//for _, bidWaiting := range bidWaitings {
			//
			//}

		}
	}

}
