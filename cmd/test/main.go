package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"upbit-api/config"
	"upbit-api/internal/api/orders"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	//setTicker := time.NewTicker(resetTime.Sub(now))

	list := orders.GetTodayDoneList()

	l := *list

	for _, value := range l {
		fmt.Println(value)
	}

	//orders.Cancel("734c8058-40b4-41bf-b7e8-6c4893f41c3c")

	<-stopChan

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
}
