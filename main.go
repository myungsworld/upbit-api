package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"upbit-api/config"
	"upbit-api/internal/api/accounts"
)

func main() {

	// 코인 매수
	//coin := models.Market("KRW-BTC")
	//coin.BidMarketPrice("5000")

	// 전체 계좌 리스트 가져오기
	fmt.Println(accounts.GetAvailableKRW())

	// 업비트 웹소켓 연결
	//connect.Socket()

	//fmt.Println(len(config.Coins))

	//r := gin.Default()
	//rAPI := r.Group("/api")
	//
	//rAPI.GET("/accounts", api.Accounts)
	//rAPI.GET("/마켓코드조회", api.AllMarketCodes)
	//r.Run()

}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	//전역변수 초기화
	config.Init()
}
