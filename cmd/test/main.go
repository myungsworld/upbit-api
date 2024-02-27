package main

import (
	"encoding/json"
	"fmt"
	"upbit-api/config"
	"upbit-api/internal/connect"
	"upbit-api/internal/models"
)

func main() {

	// 전날대비 퍼센트 비교
	conn := connect.Socket(config.Ticker)
	//tickers := make(map[string]models.Ticker, 0)
	//arrTickers := make([]models.Ticker, 0)
	//
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		ticker := models.Ticker{}
		if err := json.Unmarshal(message, &ticker); err != nil {
			panic(err)
		}

		fmt.Println(ticker)
	}

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
}
