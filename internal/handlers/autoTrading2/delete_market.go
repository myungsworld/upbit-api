package autoTrading2

import (
	"time"
	"upbit-api/internal/api/orders"
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

	//deleteTicker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-deleteTicker.C:

			waitList := orders.WaitList()
			wl := *waitList

			// TODO : 이거 확안
			year, month, day := time.Now().Add(-9 * time.Hour).Date()

			for _, value := range wl {
				// 매수체결 대기일경우
				if value.OrdType == "limit" && value.Side == "bid" {
					// 그게 금일 대기 걸어놓은 경우
					if year == value.CreatedAt.Year() && month == value.CreatedAt.Month() && day == value.CreatedAt.Day() {
						orders.Cancel(value.Uuid)
					}

				}
			}

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
