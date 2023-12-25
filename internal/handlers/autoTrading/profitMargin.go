package autoTrading

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"upbit-api/config"
	accounts2 "upbit-api/internal/api/accounts"
	"upbit-api/internal/connect"
	"upbit-api/internal/gmail"
	"upbit-api/internal/models"
)

// 첫번쨰 값은 내 평균 구매 가격 , 두번쩨 값은 가장 높은 하락률 , 세번째 값은 가장 높은 이익률
type Margin struct {
	code       string
	avgPrice   float64
	highMargin float64
	lowMargin  float64
}

func MarginMonitoring() {

	profitMarginTicker := time.NewTicker(time.Second)

	for {
		select {
		case <-profitMarginTicker.C:
			MarginMonitoringHandler()
		}
	}
}

func MarginMonitoringHandler() {

	accounts := accounts2.Get()

	marginTable := make(map[string]Margin)

	var markets []string
	for _, account := range accounts {

		if account.Currency == "KRW" || account.Currency == "APENFT" {
			continue
		}
		markets = append(markets, fmt.Sprintf("%s-%s", account.UnitCurrency, account.Currency))

		myAvgPrice, _ := strconv.ParseFloat(account.AvgBuyPrice, 64)
		marginTable[fmt.Sprintf("%s-%s", account.UnitCurrency, account.Currency)] = Margin{
			code:     fmt.Sprintf("%s-%s", account.UnitCurrency, account.Currency),
			avgPrice: myAvgPrice,
		}
	}

	conn := connect.Socket(config.Ticker, markets...)

	marginTicker := setTimer()

	for {
		select {

		case <-marginTicker.C:
			gmail.Send("보유한 코인 최대,최저 수익률", marginBody(marginTable))
			marginTicker = setTimer()

		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				if strings.Contains(err.Error(), "connection reset by peer") {
					log.Print("Profit Margin err :", err.Error())
					continue
				}
				break
			}

			ticker := models.Ticker{}

			if err := json.Unmarshal(message, &ticker); err != nil {
				panic(err)
			}

			margin := marginTable[ticker.Code]

			// 손익 퍼센트 계산
			fluctuationRate := (ticker.TradePrice/margin.avgPrice)*100 - 100

			if fluctuationRate > 0 {
				if margin.highMargin < fluctuationRate {
					marginTable[ticker.Code] = Margin{
						code:       margin.code,
						avgPrice:   margin.avgPrice,
						highMargin: fluctuationRate,
						lowMargin:  margin.lowMargin,
					}
				}
			} else {
				if margin.lowMargin > fluctuationRate {
					marginTable[ticker.Code] = Margin{
						code:       margin.code,
						avgPrice:   margin.avgPrice,
						highMargin: margin.highMargin,
						lowMargin:  fluctuationRate,
					}
				}
			}

		}
	}

}

func setTimer() *time.Ticker {
	// 현재 시간 획득
	now := time.Now()

	startTime := now.Truncate(24 * time.Hour).Add(time.Hour * 22)

	if now.UnixNano() > startTime.UnixNano() {
		startTime = startTime.Add(24 * time.Hour)
	}

	// 지연 시간 계산
	duration := startTime.Sub(now)

	// 타이머 설정
	ticker := time.NewTicker(duration)

	return ticker
}

func marginBody(marginTable map[string]Margin) string {

	var result string

	for _, margin := range marginTable {

		result += fmt.Sprintf("%s 최고점 : %0.2f  최저점 : %02.f\n", margin.code, margin.highMargin, margin.lowMargin)
	}

	return result
}
