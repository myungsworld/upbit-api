package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"sort"
	"time"
	"upbit-api/config"
	"upbit-api/internal/connect"
	"upbit-api/internal/models"
)

// 오전 8시 59분 ( UTC 기준 자정, 이전에 가장 많이 하락한 코인 구입 )
func main() {

	log.Println("오전 8시 59분 ( UTC 기준 자정, 이전에 가장 많이 하락한 코인 구입 ) 프로그램 시작")

	// 시작 시간 설정
	timer := startAt()

	for {
		<-timer.C
		// 모든 코인 조회후 가장 하락률 높은 코인 기준으로 리스트 가져오기
		getFallenCoins()

		time.Sleep(time.Second)
		fmt.Println("끝!")
		// 다음 날 오전 8시 59분으로 타이머 재설정
		timeReset(timer)
	}

}

func startAt() *time.Timer {
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 15, 52, 00, 0, now.Location())
	// 오전 8시 59분일때 실행
	duration := startTime.Sub(time.Now())
	// 이미 지난 경우 다음날 8시 59분 실행
	if duration <= 0 {
		startTime = startTime.Add(24 * time.Hour)
		duration = startTime.Sub(time.Now())
	}
	minutes := (int)(duration.Minutes()) % 60
	seconds := (int)(duration.Seconds()) % 60

	fmt.Println(fmt.Sprintf("첫 실행시간까지 남은 시간 : %.0f시간 %d분 %d초", duration.Hours(), minutes, seconds))

	return time.NewTimer(duration)
}

func timeReset(timer *time.Timer) {
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 59, 00, 0, now.Location())
	// 오전 8시 59분일때 실행
	duration := startTime.Sub(time.Now())
	// 이미 지난 경우 다음날 8시 59분 실행
	if duration <= 0 {
		startTime = startTime.Add(24 * time.Hour)
		duration = startTime.Sub(time.Now())
	}

	minutes := (int)(duration.Minutes()) % 60
	seconds := (int)(duration.Seconds()) % 60

	fmt.Println(fmt.Sprintf("다음 실행시간까지 남은 시간 : %.0f시간 %d분 %d초", duration.Hours(), minutes, seconds))

	timer.Reset(duration)
}

func init() {
	// .env 키 가져오기
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	config.Init()

}

// 하락률 제일큰 코인들부터 가져오기
func getFallenCoins() []models.Ticker {

	conn := connect.Socket()

	tickers := make(map[string]models.Ticker, 0)

	for {

		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("소켓 캇뜨")
			break
		}

		ticker := models.Ticker{}
		//fmt.Printf("Received message: %s\n", message)
		if err := json.Unmarshal(message, &ticker); err != nil {
			panic(err)
		}

		_, ok := tickers[ticker.Code]
		if !ok {
			tickers[ticker.Code] = ticker
		}

		if len(tickers) == len(config.Markets) {
			conn.Close()
		}

	}

	var minusResult []models.Ticker

	for _, ticker := range tickers {
		if ticker.Change == "FALL" {
			minusResult = append(minusResult, ticker)
		}

	}

	sort.Slice(minusResult, func(i, j int) bool {
		return minusResult[i].SignedChangeRate < minusResult[i].SignedChangeRate
	})

	for i := 0; i < len(minusResult); i++ {
		fmt.Println(i, minusResult[i].Code, fmt.Sprintf("전일대비 등락률 : %f", minusResult[i].SignedChangeRate))
	}

	return minusResult

}
