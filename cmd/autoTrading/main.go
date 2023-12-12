package main

import (
	"encoding/json"
	"fmt"
	"upbit-api/config"
	"upbit-api/internal/connect"
	"upbit-api/internal/models"
)

var bidAskAmount int64
var transactionCnt int64

func main() {

	// 일단 실시간 거래대금을 가지고 올수 있는지 확인
	conn := connect.Socket(config.Trade)

	// stream index 에 market 기입
	chanZRX := make(chan models.Trade)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			trade := models.Trade{}
			if err = json.Unmarshal(message, &trade); err != nil {
				panic(err)
			}

			if trade.Code == "KRW-AVAX" {
				if trade.AskBid == "BID" {
					bidAskAmount += int64(trade.TradePrice * trade.TradeVolume)
				} else {
					bidAskAmount -= int64(trade.TradePrice * trade.TradeVolume)
				}
				transactionCnt++
				fmt.Println(trade, " 거래량 현황:", bidAskAmount, " 거래수:", transactionCnt)
			}
		}
	}()

	for {
		select {
		case ZRX := <-chanZRX:
			fmt.Println(ZRX)
		default:
		}

	}

}

func init() {
	config.Init()

}

func autoHandler(trade models.Trade) {
	fmt.Println(trade)
}
