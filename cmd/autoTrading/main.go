package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/connect"
	"upbit-api/internal/constants"
	"upbit-api/internal/handlers/autoTrading"
	"upbit-api/internal/models"
)

// 기본 원리
// 1. 매시간 마다 전날 대비 많이 내린 코인 6천원 매수
// 2. 구매한 코인중 -10퍼가 넘어갈시 6천원 다시 매수 , 평균 7퍼 오를시 전체 매도
// 3. 실시간으로 핫한 코인 알람
// 폭락에 최소화하는 방법 분할매도
func main() {

	// TODO : 전부 setTicker() 로 time 분기
	bidEveryHourTicker := setTicker()
	averagingDownTicker := time.NewTicker(time.Minute)
	//realTimeAlarmTicker := time.NewTicker(time.Minute)
	//detectNewCoinTicker := time.NewTicker(time.Second)
	for {
		select {
		// 매시간(기준 58분) 마다 전날 대비 많이 가장 많이 내린 코인 매수 ( -5퍼가 넘지 않으면 매수하지 않음 )
		case <-bidEveryHourTicker.C:
			Market, signedRate := getBiggestFallenCoin()
			if signedRate < -5 {
				coin := orders.Market(Market)
				coin.BidMarketPrice(constants.AutoTradingBidPrice)
				bidEveryHourTicker = setTicker()
			}
		// 	구매한 코인중 -9퍼가 넘어갈시 다시 매수 , 수익률이 7퍼가 넘을시 해당 코인 전체 매도
		case <-averagingDownTicker.C:
			autoTrading.AveragingDown()
			averagingDownTicker = time.NewTicker(time.Minute * 30)
			// 실시간으로 핫한 코인 메일 받기 ( 5분이내 등락률 10퍼 상승 or 하락 )
			//case <-realTimeAlarmTicker.C:
			//	autoTrading.RealTimeAlarm()

			// 신규코인 상장시 초단위 소켓 탐지 및 매수 후 메일 발송 및 수익률 30퍼 넘을시 매도 및 메일 발송
			//case <-detectNewCoinTicker.C:
			//	autoTrading.DetectNewCoin()
		}
	}

}

func setTicker() *time.Ticker {

	now := time.Now()

	startTime := now.Truncate(time.Hour).Add(58 * time.Minute)

	if now.UnixNano() > startTime.UnixNano() {
		startTime = startTime.Add(time.Hour)
	}

	// 9시 58분은 너무 변동이 너무 심해 제외
	if startTime.Hour() == 9 {
		startTime = startTime.Add(time.Hour)
	}

	duration := startTime.Sub(now)

	ticker := time.NewTicker(duration)

	return ticker
}

func getBiggestFallenCoin() (string, float64) {
	conn := connect.Socket(config.Ticker)
	tickers := make(map[string]models.Ticker, 0)
	arrTickers := make([]models.Ticker, 0)
	for {

		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("소켓 닫기")
			break
		}

		ticker := models.Ticker{}
		if err := json.Unmarshal(message, &ticker); err != nil {
			panic(err)
		}

		_, ok := tickers[ticker.Code]
		if !ok {

			if ticker.Code == "KRW-BTT" || ticker.Code == "KRW-SHIB" || ticker.Code == "KRW-XEC" || ticker.Code == "KRW-CTC" || ticker.Code == "KRW-ASTR" {
			} else {
				tickers[ticker.Code] = ticker
				arrTickers = append(arrTickers, ticker)
			}
		}

		if len(tickers) == len(config.Markets) {
			conn.Close()
		}

	}

	sort.Slice(arrTickers, func(i, j int) bool {
		return arrTickers[i].SignedChangeRate < arrTickers[j].SignedChangeRate
	})

	return arrTickers[0].Code, arrTickers[0].SignedChangeRate
}

func init() {
	config.Init()

}
