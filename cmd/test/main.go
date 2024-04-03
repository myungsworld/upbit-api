package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"upbit-api/config"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/datastore"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	coin := orders.Market("KRW-MNT")
	uuid := coin.AskMarketPrice("28.95193978")

	order := orders.Get(uuid)
	fmt.Println("-----")
	fmt.Println(order)
	fmt.Println("-----")
	fmt.Println(order.Trades[0].Funds)

	var integerFund int

	if strings.Contains(order.Trades[0].Funds, ".") {
		fmt.Println(". 포함")
		fund := strings.Split(order.Trades[0].Funds, ".")
		integerFund, _ = strconv.Atoi(fund[0])
	} else {
		fmt.Println(". 미포함")
		integerFund, _ = strconv.Atoi(order.Trades[0].Funds)

	}

	fmt.Println(integerFund)

	<-stopChan

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
	datastore.ConnectDB()
}
