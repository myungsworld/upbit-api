package main

import (
	"os"
	"os/signal"
	"syscall"
	"upbit-api/config"
)

// 비트 시세 알람 프로그램
func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	//conn := connect.Socket(config.Ticker)

	// 비트코인이 1퍼 내리면 알람

	// 내가 가진 코인중 양전환 됐을떄 알람

	<-stopChan
}

func init() {
	config.Init()
}
