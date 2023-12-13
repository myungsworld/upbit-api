package config

import (
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	UpbitWebSocketURL = "wss://api.upbit.com/websocket/v1"

	// 소켓 데이터 요청 형식
	Ticker    = "ticker"
	Trade     = "trade"
	OrderBook = "orderbook"
)

var AccessKey string
var SecretKey string

// Markets
// 비트토렌트와 시바이누 코인은 가격 변동율이 너무 커서 제외
var Markets []string

var Coins []tradableMarket

type tradableMarket struct {
	Market        string `json:"market"`
	KoreanName    string `json:"korean_name"`
	EnglishName   string `json:"english_name"`
	MarketWarning string `json:"market_warning"`
}

func Init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	AccessKey = os.Getenv("AccessKey")
	SecretKey = os.Getenv("SecretKey")

	// 주문가능 코인들 가져오기
	getAvailableCoins()

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
