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
	"upbit-api/internal/models"
)

func AveragingDown() {

	averagingDownTicker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-averagingDownTicker.C:
			averagingDownHandler()
			averagingDownTicker = time.NewTicker(time.Minute * 30)
		}
	}
}

// 구매한 코인중 -9퍼가 넘어갈시 다시 매수 , 수익률이 7퍼가 넘을시 해당 코인 전체 매도
func averagingDownHandler() {
	accounts := accounts2.Get()

	tickers := getCurrentTickerMappingAccounts(accounts)

	for _, account := range accounts {

		ticker := tickers[fmt.Sprintf("%s-%s", account.UnitCurrency, account.Currency)]

		myAvgPrice, _ := strconv.ParseFloat(account.AvgBuyPrice, 64)

		// 손익 퍼센트 계산
		fluctuationRate := (ticker.TradePrice/myAvgPrice)*100 - 100

		// -9 퍼가 넘을시 다시 매수
		if fluctuationRate < -9 {
			coin := orders.Market(ticker.Code)
			coin.BidMarketPrice(constants.AutoTradingBidPrice)
			log.Print(coin, fmt.Sprintf(" %0.2f%% 매수", fluctuationRate))
		}

		// 15 퍼 이상 수익률 나면 매도
		// TODO : 아직 이걸 못정하겠음
		if fluctuationRate > 15 {
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
		}

	}

	return tickers
}
