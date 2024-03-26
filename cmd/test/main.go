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

	//setTicker := time.NewTicker(resetTime.Sub(now))

	datastore.ConnectDB()

	//TODO : 9시 전의 매수 체결 대기 걸어놓고 테스트 해보기

	//// 그날의 모든 매수 체결 대기 가져오기
	//currentTime := time.Now().UTC()
	//startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)
	//endOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, time.UTC)
	//bidWaitings := make([]models.BidWaiting, 0)
	//if err := datastore.DB.Model(&models.BidWaiting{}).
	//	Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).
	//	Find(&bidWaitings).Error; err != nil {
	//	panic(err)
	//}
	//
	//// 매수 체결 대기 제거
	//for _, bidWaiting := range bidWaitings {
	//	orders.Cancel(bidWaiting.Uuid)
	//	if err := datastore.DB.Delete(&bidWaiting).Error; err != nil {
	//		panic(err)
	//	}
	//
	//}

	orders.Get("1edd38db-10a1-4f28-ac6a-efc4db0bb610")

	<-stopChan

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
}
