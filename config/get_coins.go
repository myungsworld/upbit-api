package config

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var Coins []tradableMarket

type tradableMarket struct {
	Market        string `json:"market"`
	KoreanName    string `json:"korean_name"`
	EnglishName   string `json:"english_name"`
	MarketWarning string `json:"market_warning"`
}

// getAvailableCoins
// https://api.upbit.com/v1/market/all
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

	err = json.Unmarshal(body, &Coins)
	if err != nil {
		panic(err)
	}

	// 원화 마켓만 리스트 업데이트
	for _, coin := range Coins {

		if coin.Market[0:3] == "KRW" {
			// 아래 코인 가격이 너무 낮아 변동률이 커서 제외
			if coin.Market == "KRW-BTT" || coin.Market == "KRW-SHIB" {
			} else {
				Markets = append(Markets, coin.Market)
			}

		}
	}

}
