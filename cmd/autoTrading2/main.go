package main

import (
	"os"
	"os/signal"
	"syscall"
	"upbit-api/config"
	"upbit-api/internal/datastore"
	"upbit-api/internal/handlers/autoTrading2"
)

// 코인 스윙 자동매매
// 지나간 n일의 저점,종가,고가의 평균을 가져와 종가의 평균과 현재가의 편차가 적을때
// 저점의 평균에서 매수대기 이후 매수가 되면 고점의 평균 -1퍼에서 매도
// 하루동안 모니터링후 매수대기만 걸린 데이터들 주문리스트에서 삭제후 초기화 및 회귀
func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// TODO: 금일 기준 데이터가 데이터베이스에 들어가있지 않은 경우 기입 ( 프로그램이 멈추고 난후 매수나 매도가 일어난 경우 )
	autoTrading2.SetDatabase()

	// 1일마다 리셋 ( 한국시각 9시 1초 )
	// 3일 저점,종가,고가 평균 연산 후 상태값 저장
	go autoTrading2.Reset()

	// 초기화된 데이터 지정가 매수 체결 대기 걸기
	go autoTrading2.LimitOrder()

	// 매수가 되었다면 매도 대기 걸기
	go autoTrading2.AskOrder()

	// 매도가 되었다면 데이터베이스 저장
	go autoTrading2.AskCheck()

	// 3시 , 6시 폭락한 코인 매수
	//go autoTrading2.()

	// 매일 8시 55분 매수체결대기와 매도체결대기가 계속 걸려 있을시 그날의 매수체결 대기 삭제 매도 체결대기는 매도가 될때까지 상태 유지
	go autoTrading2.DeleteWaitMarket()

	// (6시간마다 실행) 매도가 해당 날이 아닌 다른 날에 되었을 경우 데이터베이스 업데이트
	go autoTrading2.UpdateDB()

	<-stopChan

}

func init() {
	config.Init()
	datastore.ConnectDB()
}
