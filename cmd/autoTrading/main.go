package main

import (
	"fmt"
	"log"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/constants"
	"upbit-api/internal/handlers/autoTrading"
)

// 기본 매매 방식
// 1. 매시간 마다 전날 대비 많이 내린 코인 6천원 매수
// 2. 구매한 코인중 -9퍼가 넘어갈시 6천원 다시 매수 , 평균 7퍼 오를시 전체 매도
// 3. 실시간으로 핫한 코인 알람
// 폭락에 최소화하는 방법 분할매도
func main() {

	bidEveryHourTicker := autoTrading.SetTickerForBidEveryHour()
	averagingDownTicker := time.NewTicker(time.Second)
	detectNewCoinTicker := time.NewTicker(time.Second)
	//realTimeAlarmTicker := time.NewTicker(time.Minute)
	for {
		select {
		// 매시간(기준 58분) 마다 전날 대비 많이 가장 많이 내린 코인 매수 ( -5퍼가 넘지 않으면 매수하지 않음 )
		case <-bidEveryHourTicker.C:
			market, signedRate := autoTrading.GetBiggestFallenCoin()
			if signedRate < -5 {
				coin := orders.Market(market)
				coin.BidMarketPrice(constants.AutoTradingBidPrice)
				log.Print(market, fmt.Sprintf(" %0.2f%% 하락 , 매수 %s원", signedRate, constants.AutoTradingBidPrice))
			}
			bidEveryHourTicker = autoTrading.SetTickerForBidEveryHour()

		// 	구매한 코인중 -9퍼가 넘어갈시 다시 매수 , 수익률이 7퍼가 넘을시 해당 코인 전체 매도
		case <-averagingDownTicker.C:
			autoTrading.AveragingDown()
			averagingDownTicker = time.NewTicker(time.Second)

		//신규코인 상장시 탐지 및 매수 후 메일 발송 및 수익률 30퍼 넘을시 매도 및 메일 발송
		case <-detectNewCoinTicker.C:
			autoTrading.DetectNewCoin()

			// 급등 급락 탐지

		}
	}

}

func init() {
	config.Init()

}
