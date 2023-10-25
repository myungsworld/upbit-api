package main

import (
	"fmt"
	"github.com/goccy/go-json"
	"upbit-api/config"
	"upbit-api/internal/connect"
	"upbit-api/internal/models"
)

func main() {

	// 일단 실시간 거래대금을 가지고 올수 있는지 확인
	conn := connect.Socket(config.OrderBook)

	for {

		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		orderBook := models.OrderBook{}
		if err = json.Unmarshal(message, &orderBook); err != nil {
			panic(err)
		}

		if orderBook.Code == "KRW-BTC" {
			fmt.Println(orderBook)
		}

	}

}

func init() {
	config.Init()
}
