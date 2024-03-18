package autoTrading2

import (
	"time"
	"upbit-api/internal/api/orders"
)

func DeleteWaitMarket() {

	// 매일 8시 55분 초기화
	//now := time.Now()
	//resetTime := now.Truncate(24 * time.Hour).Add(time.Hour * 24).Add(-5 * time.Minute)
	//deleteTicker := time.NewTicker(resetTime.Sub(now))

	deleteTicker := time.NewTicker(time.Second * 60)

	for {
		select {
		case <-deleteTicker.C:

			waitList := orders.WaitList()
			wl := *waitList

			year, month, day := time.Now().Date()

			for _, value := range wl {
				// 매수체결 대기일경우
				if value.OrdType == "limit" && value.Side == "bid" {
					// 그게 금일 대기 걸어놓은 경우
					if year == value.CreatedAt.Year() && month == value.CreatedAt.Month() && day == value.CreatedAt.Day() {
						orders.Cancel(value.Uuid)
					}

				}
			}

			//now := time.Now()
			//resetTime := now.Truncate(24 * time.Hour).Add(time.Hour * 24).Add(-5 * time.Minute)
			//deleteTicker := time.NewTicker(resetTime.Sub(now))

		}
	}

}
