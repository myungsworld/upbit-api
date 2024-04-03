package middlewares

import (
	"fmt"
	"time"
)

// SetTimerEveryHourByMinute 매시간마다 돌아가는 타이머
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

// SetTimerEvery6Hour 6시간마다 돌아가는 타이머 (정오 기준)
func SetTimerEvery6Hour() *time.Ticker {

	nowHour := time.Now().UTC().Hour()

	fmt.Println(nowHour)

	now := time.Now()
	resetTime := time.Now().UTC().Truncate(time.Hour * 24)

	// 00시, 06시, 12시, 18시 기준으로 구매
	switch {
	case 0 <= nowHour && nowHour < 6:
		resetTime = resetTime.Add(time.Hour * 6)
	case 6 <= nowHour && nowHour < 12:
		resetTime = resetTime.Add(time.Hour * 12)
	case 12 <= nowHour && nowHour < 18:
		resetTime = resetTime.Add(time.Hour * 18)
	case 18 <= nowHour && nowHour < 24:
		resetTime = resetTime.Add(time.Hour * 24)
	}

	ticker := time.NewTicker(resetTime.Sub(now))

	return ticker

}
