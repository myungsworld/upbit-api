package config

import (
	"os"
)

const (
	UpbitWebSocketURL = "wss://api.upbit.com/websocket/v1"
)

// Markets
// 비트토렌트와 시바이누 코인은 가격 변동율이 너무 커서 제외
var Markets []string

var AccessKey string
var SecretKey string

func Init() {

	AccessKey = os.Getenv("AccessKey")
	SecretKey = os.Getenv("SecretKey")

	// 주문가능 코인들 가져오기
	getAvailableCoins()

}
