package main

import (
	"fmt"
	"time"
	"upbit-api/config"
)

func main() {

	for i := 0; i < 20; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("혁두캔두잇")
	}

	//// 전날대비 퍼센트 비교
	//conn := connect.Socket(config.Ticker)
	////tickers := make(map[string]models.Ticker, 0)
	////arrTickers := make([]models.Ticker, 0)
	////
	//for {
	//	_, message, err := conn.ReadMessage()
	//	if err != nil {
	//		break
	//	}
	//
	//	ticker := models.Ticker{}
	//	if err := json.Unmarshal(message, &ticker); err != nil {
	//		panic(err)
	//	}
	//
	//	fmt.Println(ticker)
	//}

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
}
