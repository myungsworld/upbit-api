package main

import (
	"github.com/joho/godotenv"
	"log"
	"upbit-api/config"
	"upbit-api/internal/connect"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// 전역변수 초기화
	config.Init()

	// 업비트 웹소켓 연결
	connect.Socket()

	//fmt.Println(len(config.Coins))

	//r := gin.Default()
	//rAPI := r.Group("/api")
	//
	//rAPI.GET("/accounts", handlers.Accounts)
	//rAPI.GET("/마켓코드조회", handlers.AllMarketCodes)
	//r.Run()

}
