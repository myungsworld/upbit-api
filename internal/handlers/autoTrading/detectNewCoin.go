package autoTrading

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"upbit-api/config"
	"upbit-api/internal/gmail"
)

type NewTradableMarket struct {
	Market string `json:"market"`
}

var newCoins []NewTradableMarket

var newMarkets []string

func DetectNewCoin() {
	originMarketPool := make(map[string]int)

	for _, Market := range config.Markets {

		originMarketPool[Market] = 1
	}
	for {
		newMarkets = make([]string, 0)
		getAvailableCoins()
		// 원래는 < 이렇게 해야 새로운 마켓이 더많아서 새로운 코인이 추가되었다고 보는데
		// 그 외의 다른경우도 있을까 싶어서
		if len(config.Markets) != len(newMarkets) {

			for _, newMarket := range newMarkets {
				if _, exist := originMarketPool[newMarket]; exist {

				} else {
					// TODO : config.Markets 초기화
					// TODO : newMarket 매수 후 대기열 진입
					gmail.Send(fmt.Sprintf("신규코인 상장 :%s", newMarket), "")
					time.Sleep(time.Second * 10)
				}
			}
		}
		time.Sleep(time.Second * 1)

		fmt.Print("이전:", len(config.Markets), ",")
		fmt.Println("이후:", len(newMarkets))
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
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
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
			if coin.Market == "KRW-BTT" || coin.Market == "KRW-SHIB" || coin.Market == "KRW-XEC" {
			} else {
				newMarkets = append(newMarkets, coin.Market)
			}

		}
	}

}
