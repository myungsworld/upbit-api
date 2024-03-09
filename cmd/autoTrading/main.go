package main

import (
	"os"
	"os/signal"
	"syscall"
	"upbit-api/config"
	"upbit-api/internal/handlers/autoTrading"
)

// 기본 매매 방식
// 1. 매시간 마다 전날 대비 많이 내린 코인 6천원 매수
// 2. 구매한 코인중 -9퍼가 넘어갈시 6천원 다시 매수 , 평균 7퍼 오를시 전체 매도
// 3. 실시간으로 핫한 코인 알람
// 폭락에 최소화하는 방법 분할매도
func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// health check
	go autoTrading.HealthCheck()
	// 매시간(기준 58분) 마다 전날 대비 많이 가장 많이 내린 코인 매수 ( -5퍼가 넘지 않으면 매수하지 않음 )
	go autoTrading.BidEveryHour()
	// 구매한 코인의 수익률이 -12퍼가 넘어갈시 다시 매수 ( 6시간 마다 )
	go autoTrading.AveragingDown()
	// 구매한 코인중 수익률이 20퍼가 넘을시 해당 코인 전체 매도
	go autoTrading.AskCoin()
	// 구매한 코인중 -5퍼 손해일시 해당 코인 금액에 따라 부분 매도

	// 내가 가진 코인중 금일 많이 오른 코인 부분 매도

	// 신규 코인 상장시 탐지 및 매수 후 메일 발송 및 수익률 30퍼 넘을시 매도 및 메일 발송
	//go autoTrading.DetectNewCoin()

	<-stopChan

}

func init() {
	config.Init()
}
