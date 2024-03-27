package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"upbit-api/config"
	"upbit-api/internal/datastore"
	"upbit-api/internal/handlers/autoTrading2"
	"upbit-api/internal/models"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	flows := make([]models.AutoTrading2, 0)

	if err := datastore.DB.Model(&models.AutoTrading2{}).
		Where("w_deleted_at is null").
		Find(&flows).Error; err != nil {
		panic(err)
	}

	for _, flow := range flows {
		lowHighGap := flow.HighTradeGap - flow.LowTradeGap - 1
		askPercent := lowHighGap/100 + 1
		fmt.Println("--", flow.Ticker, "--")
		bidFloat, _ := strconv.ParseFloat(flow.BidPrice, 64)
		askFloat := bidFloat * askPercent
		fmt.Println(flow.BidPrice, flow.HighTradeGap, flow.LowTradeGap)
		fmt.Println(lowHighGap/100 + 1)

		fmt.Println(flow.BidPrice)
		fmt.Println(bidFloat)

		fmt.Println(bidFloat * askPercent)

		fmt.Println(autoTrading2.SetAskPrice(askFloat))

		fmt.Println("--")

		askPrice := autoTrading2.SetAskPrice(askFloat)

		if err := datastore.DB.Model(&flow).
			Update("ask_price", askPrice).Error; err != nil {
			panic(err)
		}

	}

	<-stopChan

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
	datastore.ConnectDB()
}
