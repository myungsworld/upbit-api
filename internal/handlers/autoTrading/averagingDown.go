package autoTrading

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
	"upbit-api/config"
	accounts2 "upbit-api/internal/api/accounts"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/connect"
	"upbit-api/internal/constants"
	"upbit-api/internal/middlewares"
	"upbit-api/internal/models"
)

func AskCoin() {
	askCoinTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-askCoinTicker.C:
			askCoinHandler()
			askCoinTicker = time.NewTicker(time.Second * 2)
		}
	}
}

func AveragingDown() {

	averagingDownTicker := middlewares.SetTimerEvery6Hour()

	for {
		select {
		case <-averagingDownTicker.C:
			averagingDownHandler()
			averagingDownTicker = middlewares.SetTimerEvery6Hour()
		}
	}
}

// 구매한 코인중 -12퍼가 넘어갈시 다시 매수
func averagingDownHandler() {
	accounts := accounts2.Get()

	tickers := getCurrentTickerMappingAccounts(accounts)

	for _, account := range accounts {

		ticker := tickers[fmt.Sprintf("%s-%s", account.UnitCurrency, account.Currency)]

		myAvgPrice, _ := strconv.ParseFloat(account.AvgBuyPrice, 64)

		// 손익 퍼센트 계산
		fluctuationRate := (ticker.TradePrice/myAvgPrice)*100 - 100

		// -9 퍼가 넘을시 다시 매수
		if fluctuationRate < -12 {
			coin := orders.Market(ticker.Code)
			coin.BidMarketPrice(constants.AutoTradingBidPrice)
			log.Print(coin, fmt.Sprintf(" %0.2f%% 매수", fluctuationRate))
		}

		// 15 퍼 이상 수익률 나면 매도
		//if fluctuationRate > 15 {
		//	coin := orders.Market(ticker.Code)
		//	coin.AskMarketPrice(account.Balance)
		//	log.Print(coin, fmt.Sprintf(" %0.2f%% 매도", fluctuationRate))
		//}

	}

}

func askCoinHandler() {
	accounts := accounts2.Get()

	tickers := getCurrentTickerMappingAccounts(accounts)

	for _, account := range accounts {

		ticker := tickers[fmt.Sprintf("%s-%s", account.UnitCurrency, account.Currency)]

		myAvgPrice, _ := strconv.ParseFloat(account.AvgBuyPrice, 64)

		// 손익 퍼센트 계산
		fluctuationRate := (ticker.TradePrice/myAvgPrice)*100 - 100

		// 12 퍼 이상 수익률 나면 메세지 보내기

		// 20퍼 이상 오른 경우 매도
		if fluctuationRate > 20 {
			// (비트코인 제외)
			if ticker.Code == "KRW-BTC" {
				continue
			}
			coin := orders.Market(ticker.Code)
			coin.AskMarketPrice(account.Balance)
			log.Print(coin, fmt.Sprintf(" %0.2f%% 매도", fluctuationRate))
		}

	}
}

func getCurrentTickerMappingAccounts(accounts models.Accounts) map[string]models.Ticker {

	var markets []string

	for _, account := range accounts {

		if account.Currency == "KRW" || account.Currency == "APENFT" {
			continue
		}
		markets = append(markets, fmt.Sprintf("%s-%s", account.UnitCurrency, account.Currency))
	}

	conn := connect.Socket(config.Ticker, markets...)

	tickers := make(map[string]models.Ticker)
	arrTickers := make([]models.Ticker, 0)

	for {

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("여기인가2", err)
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

		if len(tickers) == len(markets) {
			conn.Close()
			break
		}

	}

	return tickers
}
