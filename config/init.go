package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	UpbitWebSocketURL = "wss://api.upbit.com/websocket/v1"
)

var AccessKey string
var SecretKey string

func Init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	AccessKey = os.Getenv("AccessKey")
	SecretKey = os.Getenv("SecretKey")

	// 주문가능 코인들 가져오기
	getAvailableCoins()

	//

}
