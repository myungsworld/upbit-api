package autoTrading

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"upbit-api/config"
	"upbit-api/internal/api/accounts"
	"upbit-api/internal/api/orders"
	"upbit-api/internal/connect"
	"upbit-api/internal/constants"
	"upbit-api/internal/gmail"
	"upbit-api/internal/models"
)

type NewTradableMarket struct {
	Market string `json:"market"`
}

var newCoins []NewTradableMarket

var newMarkets []string

func DetectNewCoin() {

	detectNewCoinTicker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-detectNewCoinTicker.C:
			detectNewCoinHandler()
		}
	}
}

func detectNewCoinHandler() {
	originMarketPool := make(map[string]int)

	for _, Market := range config.Markets {

		originMarketPool[Market] = 1
	}
	for {
		newMarkets = make([]string, 0)
		getAvailableCoins()

		if len(newMarkets) == 0 {
			//log.Print("detect New coin 에러 발생")
			continue
		}
		// 원래는 < 이렇게 해야 새로운 마켓이 더많아서 새로운 코인이 추가되었다고 보는데
		// 그 외의 다른경우도 있을까 싶어서
		if len(config.Markets) != len(newMarkets) {

			for _, newMarket := range newMarkets {
				if _, exist := originMarketPool[newMarket]; exist {
				} else {
					// 새로운 코인 업비트 상장

					var myAvgPrice float64
					var balance string

					// 새 코인 매수
					coin := orders.Market(newMarket)
					coin.BidMarketPrice(constants.NewCoinBidPrice)

					// 매수한 코인 정보 가져오기
					accts := accounts.Get()
					for _, acct := range accts {
						if newMarket == fmt.Sprintf("%s-%s", acct.UnitCurrency, acct.Currency) {
							myAvgPrice, _ = strconv.ParseFloat(acct.AvgBuyPrice, 64)
							balance = acct.Balance
							break
						}
					}

					// 금액이 없어 주문을 못한 경우 , 새 코인 찾기
					if myAvgPrice == 0 && balance == "" {

						// 실패 메일 전송
						gmail.Send(fmt.Sprintf("신규 %s 코인 상장", newMarket), "금액 부족으로 인한 구매 실패")
						// config.Markets 초기화
						config.Markets = append(config.Markets, newMarket)
						originMarketPool[newMarket] = 1
						break
					}

					// 매수 성송 메일 전송
					gmail.Send(fmt.Sprintf("신규 %s 코인 상장", newMarket),
						fmt.Sprintf("%s\n매수금액 : %s원 \n매수평균가: %f원", newMarket, constants.NewCoinBidPrice, myAvgPrice))

					// config.Markets 초기화
					config.Markets = append(config.Markets, newMarket)
					originMarketPool[newMarket] = 1

					// 소켓 열고 30퍼 이상일시 매도 후 소켓 종료 or 1시간 지나면 소켓 종료
					bidNewPublicMarket(newMarket, balance, myAvgPrice)
					break
				}
			}

		}
		// 기존 마켓이랑 새로운 마켓의 길이가 어떻게 다른지 확인
		//fmt.Println(len(config.Markets), len(newMarkets))
		time.Sleep(time.Second * 1)
	}
}

func getAvailableCoins() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.upbit.com/v1/market/all?isDetails=true", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			//gmail.Send("detectNewCoin 에러", err.Error())
			return
		}
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			//gmail.Send("detectNewCoin 에러", err.Error())
			return
		}
		panic(err)
	}

	err = json.Unmarshal(body, &newCoins)
	if err != nil {
		panic(err)
	}

	// 원화 마켓만 리스트 업데이트
	for _, coin := range newCoins {

		if coin.Market[0:3] == "KRW" {
			// 아래 코인 가격이 너무 낮아 변동률이 커서 제외
			if config.ExceptMarkets(coin.Market) {
			} else {
				newMarkets = append(newMarkets, coin.Market)
			}

		}
	}

}

func bidNewPublicMarket(market, balance string, myAvgPrice float64) {

	conn := connect.Socket(config.Ticker, market)

	log.Printf("%s 신규코인 매도 대기열로 들어옴", market)

	timer := time.NewTimer(time.Second * 5)

	// 대기열
	for {
		// 1시간이 지나면 해당 대기열 종료
		select {
		case <-timer.C:
			timer.Stop()
			log.Printf("%s 신규코인 매도 대기열 종료", market)
			conn.Close()
			return

		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			ticker := models.Ticker{}

			if err := json.Unmarshal(message, &ticker); err != nil {
				panic(err)
			}

			fluctuationRate := (ticker.TradePrice/myAvgPrice)*100 - 100

			// 매수 평균 가격보다 30퍼 이상 올랐을시 전체 매도
			if fluctuationRate > 30 {
				coin := orders.Market(market)
				coin.AskMarketPrice(balance)

				conn.Close()
			}

		}
	}

}
