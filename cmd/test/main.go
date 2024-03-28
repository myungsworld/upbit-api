package main

import (
	"os"
	"os/signal"
	"syscall"
	"upbit-api/config"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/datastore"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	orders.Get("4ab91eae-fbf6-4e05-b577-e54a5a5cd538")

	<-stopChan

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
	datastore.ConnectDB()
}
