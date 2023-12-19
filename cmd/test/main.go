package main

import (
	"upbit-api/config"
	"upbit-api/internal/api/candle"
)

func main() {

	//acc := accounts.Get()
	//fmt.Println(acc)

	coin := candle.Market("KRW-BTC")
	coin.Min()
}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
}
