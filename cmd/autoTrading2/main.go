package main

import (
	"os"
	"os/signal"
	"syscall"
	"upbit-api/config"
	"upbit-api/internal/handlers/autoTrading2"
)

// 코인 스윙 자동매매
// 지나간 n일의 저점,종가,고가의 평균을 가져와 종가의 평균과 현재가의 편차가 적을때
// 저점의 평균에서 매수대기 이후 매수가 되면 고점의 평균 -1퍼에서 매도
// 하루동안 모니터링후 매수대기만 걸린 데이터들 주문리스트에서 삭제후 초기화 및 회귀
func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 1일마다 리셋 ( 한국시각 9시 1초 )
	// 3일 저점,종가,고가 평균 연산 후 상태값 저장
	go autoTrading2.Reset()

	// 초기화된 데이터 지정가 매수 체결 대기 걸기
	go autoTrading2.LimitOrder()

	// 매일 8시 55분 매수체결 대기가 계속 걸려 있을시 그날의 매수체결 대기 삭제
	// TODO : 이거 테스트 해봐야함
	go autoTrading2.DeleteWaitMarket()

	// 고점의 종가 -1%에서 매도
	//go autoTrading2.Handler()

	<-stopChan

}

func init() {
	config.Init()
}
