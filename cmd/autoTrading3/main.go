package main

import (
	"os"
	"os/signal"
	"syscall"
	"upbit-api/config"
	"upbit-api/internal/datastore"
	"upbit-api/internal/handlers/autoTrading3"
)

// autoTrading2 에서 매수기준만 변경된 프로그램
func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 1일마다 리셋 ( 한국시각 9시 1초 )
	// 3일 저점,종가,고가 평균 연산 후 상태값 저장
	go autoTrading3.Reset()

	// 초기화된 데이터 지정가 매수 체결 대기 걸기 및 종가 평균 대비 현재가 계속 업데이트
	go autoTrading3.LimitOrder()

	// 매수가 되었다면 매도 대기 걸기
	//go autoTrading3.AskOrder()

	// 매도가 되었다면 데이터베이스 저장
	//go autoTrading3.AskCheck()

	//TODO: 3시 , 12시 현재가가 저점의 평균보다 낮은 코인 매수
	go autoTrading3.BidBigShort()

	// 매일 8시 55분 매수체결대기와 매도체결대기가 계속 걸려 있을시 그날의 매수체결 대기 삭제 매도 체결대기는 매도가 될때까지 상태 유지
	//go autoTrading3.DeleteWaitMarket()

	// (6시간마다 실행) 매도가 해당 날이 아닌 다른 날에 되었을 경우 데이터베이스 업데이트
	//go autoTrading3.UpdateDB()

	<-stopChan

}

func init() {
	config.Init()
	datastore.ConnectDB()
}
