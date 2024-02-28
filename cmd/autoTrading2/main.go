package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/candle"
)

type info struct {
	LowAverage  float64
	HighAverage float64
}

// 이전 데이터 가져오는 기준일
const count = 4

var previousMarketInfo map[string]info

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 상태값 담을 map
	previousMarketInfo = make(map[string]info)

	// 1일마다 리셋 ( 한국시각 9시 )
	go reset()

	//oneSecTicker := time.NewTicker(time.Second)
	//
	//go func() {
	//	for {
	//		select {
	//		case <-oneSecTicker.C:
	//			fmt.Println(previousMarketInfo)
	//		}
	//	}
	//}()

	// 저점의 평균으로 왔을때 매수

	// 고점과 저점의 퍼센트 차이 구하고 그 반 퍼센트의 수익이 난다면 매도

	<-stopChan

}

func reset() {

	setTicker := time.NewTicker(time.Second)

	//
	for {
		select {
		case <-setTicker.C:

			log.Println("데이터 초기화 시작")

			// 데이터 초기화 실행
			// 상태값 담을 map
			previousMarketInfo = make(map[string]info)

			// 지난 3일의 저점과 고점의 평균을 구함 ( 하루동안 상태값으로 남겨놓음 )

			market := candle.Market("KRW-BTC")

			market.Day(count)

			//for market := range config.Markets {
			//
			//}

			// 매일 9시 티커 설정
			now := time.Now()

			resetTime := now.Truncate(24 * time.Hour).Add(time.Hour * 24).Add(time.Second)

			setTicker = time.NewTicker(resetTime.Sub(now))

			log.Println("데이터 초기화 끝")

		}
	}

}

func init() {
	config.Init()
}
