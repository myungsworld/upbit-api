package autoTrading

import (
	"encoding/json"
	"log"
	"sort"
	"time"
	"upbit-api/config"
	"upbit-api/internal/connect"
	"upbit-api/internal/models"
)

func GetBiggestFallenCoin() (string, float64) {
	conn := connect.Socket(config.Ticker)

	tickers := make(map[string]models.Ticker, 0)
	arrTickers := make([]models.Ticker, 0)
	for {

		_, message, err := conn.ReadMessage()

		if err != nil {
			break
		}

		ticker := models.Ticker{}
		if err := json.Unmarshal(message, &ticker); err != nil {
			panic(err)
		}

		_, ok := tickers[ticker.Code]
		if !ok {

			if config.ExceptMarkets(ticker.Code) {
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

	return arrTickers[0].Code, arrTickers[0].SignedChangeRate * 100
}

func SetTickerForBidEveryHour() *time.Ticker {

	now := time.Now()

	startTime := now.Truncate(time.Hour).Add(58 * time.Minute)

	if now.UnixNano() > startTime.UnixNano() {
		startTime = startTime.Add(time.Hour)
	}

	duration := startTime.Sub(now)

	log.Print(duration)

	ticker := time.NewTicker(duration)

	return ticker
}
