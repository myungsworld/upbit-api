package main

import (
	"upbit-api/config"
	"upbit-api/internal/api/orders"
)

func main() {

	coin := orders.Market("KRW-SEI")
	coin.BidMarketPrice("6000")
}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
}
