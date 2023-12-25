package main

import (
	"os"
	"os/signal"
	"syscall"
	"upbit-api/config"
	"upbit-api/internal/handlers/autoTrading"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	//log.Print("고난이 대수랴 기꺼이 걸으리 수라의 길")
	// 코인
	{
		// health check
		go autoTrading.HealthCheck()
		// 매시간(기준 58분) 마다 전날 대비 많이 가장 많이 내린 코인 매수 ( -5퍼가 넘지 않으면 매수하지 않음 )
		go autoTrading.BidEveryHour()
		// 구매한 코인중 -9퍼가 넘어갈시 다시 매수 , 수익률이 7퍼가 넘을시 해당 코인 전체 매도
		go autoTrading.AveragingDown()
		// 신규 코인 상장시 탐지 및 매수 후 메일 발송 및 수익률 30퍼 넘을시 매도 및 메일 발송
		go autoTrading.DetectNewCoin()
		// 구매한 코인의 역대 최고수익률 ,최저점 기록 -> 정오에 메일 발송
		//go autoTrading.MarginMonitoring()
		// DAO 자동 공모
		// go autoTrading.Dao()
		// 급등 급락 코인 바이패스
		// go autoTrading.ByPass()
	}
	// 주식
	{
		// 한국 전자공시 ( Dart ) 연결재무제표 크롤링
		// go dart.FinancialStatKr()
		// 미국 전자공시 ( EDGAR ) 연결재무제표 크롤링
		// go dart.FinancialStatUS()
		// YahooFinance 기사 크롤링 -> 감정분석 모델 -> 긍정반응일시 해당 Stock 매수
		// go autoTrading.Stock()

	}
	// 부동산
	{
		// 아이디어 가온나 이말이다 궤이쉐이야
	}
	// 원하는 시간대 기차표 자동으로 구입
	// go korail.BidTicket()
	<-stopChan

}

func init() {
	config.Init()
}
