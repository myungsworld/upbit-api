package middlewares

import (
	"time"
)

// 매시간마다 돌아가는 타이머
func SetTimerEveryHourByMinute(min int) *time.Ticker {

	now := time.Now()

	setDuration := time.Duration(min) * time.Minute

	startTime := now.Truncate(time.Hour).Add(setDuration)

	if now.UnixNano() > startTime.UnixNano() {
		startTime = startTime.Add(time.Hour)
	}

	duration := startTime.Sub(now)
	ticker := time.NewTicker(duration)

	return ticker

}
