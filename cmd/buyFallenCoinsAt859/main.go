package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/accounts"
	"upbit-api/internal/connect"
	"upbit-api/internal/models"
)

// 오전 8시 59분 ( UTC 기준 자정 기준 가장 많이 하락한 코인 구입 )
func main() {

	log.Println("( UTC 기준 자정 기준 가장 많이 하락한 코인 구입 ) 프로그램 시작")

	// 시작 시간 설정
	timer := startAt()

	for {
		<-timer.C
		// 모든 코인 조회후 가장 하락률 높은 코인 기준으로 리스트 가져오기
		tickers := getFallenCoins()

		purchaseAmountStr := "6000"
		buyTickers := tickers[:1]

		purchaseAmount, _ := strconv.Atoi(purchaseAmountStr)

		// 총 현금보다 적으면 못삼
		if accounts.GetAvailableKRW() < purchaseAmount*len(buyTickers) {
			fmt.Println("현재 가지고 있는 한화:", accounts.GetAvailableKRW(), "원")
			fmt.Println("주문 가능 금액:", purchaseAmount*len(buyTickers), "원")
		} else {
			for _, ticker := range buyTickers {
				coin := models.Market(ticker.Code)
				coin.BidMarketPrice(purchaseAmountStr)
			}
		}

		// 사는거 구현해야함
		// lambda_deploy_buyFallenCoinsAt859()

		// 다음 날 오전 8시 59분으로 타이머 재설정
		timeReset(timer)
	}

}

func startAt() *time.Timer {
	startTime := time.Now().Add(time.Second)
	//startTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 59, 00, 0, now.Location())
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
	//err := godotenv.Load("../../.env")
	//if err != nil {
	//	log.Fatal("Error loading .env file:", err)
	//}

	config.Init()

}

// 하락률 제일큰 코인들부터 가져오기
func getFallenCoins() []models.Ticker {

	// 현재가 가져오는 소켓 연결
	conn := connect.Socket(config.Ticker)

	tickers := make(map[string]models.Ticker, 0)
	arrTickers := make([]models.Ticker, 0)
	for {

		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("소켓 캇뜨")
			break
		}

		ticker := models.Ticker{}
		if err := json.Unmarshal(message, &ticker); err != nil {
			panic(err)
		}

		_, ok := tickers[ticker.Code]
		if !ok {
			tickers[ticker.Code] = ticker
			arrTickers = append(arrTickers, ticker)
		}

		if len(tickers) == len(config.Markets) {
			conn.Close()
		}

	}

	sort.Slice(arrTickers, func(i, j int) bool {
		return arrTickers[i].SignedChangeRate < arrTickers[j].SignedChangeRate
	})

	//for _, arrTicker := range arrTickers {
	//	fmt.Println(arrTicker.Code, arrTicker.SignedChangeRate)
	//}

	return arrTickers

}
