package main

import (
	"fmt"
	"upbit-api/config"
	"upbit-api/internal/api/accounts"
)

func main() {

	acc := accounts.Get()
	fmt.Println(acc)

	//coin := orders.Market("KRW-APT")
	//coin.AskMarketPrice("0.49140049")
}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
}
