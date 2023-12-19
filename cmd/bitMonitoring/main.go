package main

import (
	"encoding/json"
	"fmt"
	"upbit-api/config"
	"upbit-api/internal/connect"
	"upbit-api/internal/models"
)

func main() {

	conn := connect.Socket(config.Ticker, "KRW-BTC")

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			ticker := models.Ticker{}

			if err = json.Unmarshal(message, &ticker); err != nil {
				panic(err)
			}

		}
	}()

	select {}
}

func init() {
	config.Init()
}
