package main

import (
	"os"
	"os/signal"
	"syscall"
	"upbit-api/config"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
}
