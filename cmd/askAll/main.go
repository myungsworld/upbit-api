package main

import (
	"fmt"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/accounts"
	"upbit-api/internal/api/orders"
)

// 비트코인 제외 가지고 있는 코인 시장가로 전부 매도
func main() {

	accts := accounts.Get()

	for _, acct := range accts {

		if acct.Currency == "BTC" || acct.Currency == "KRW" {
		} else {
			coin := orders.Market(fmt.Sprintf("KRW-%s", acct.Currency))
			coin.AskMarketPrice(acct.Balance)
			time.Sleep(time.Millisecond * 500)
		}

	}

}

func init() {
	config.Init()
}
