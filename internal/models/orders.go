package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"upbit-api/internal/middlewares"
)

const (
	// 매수
	buy = "bid"
	// 매도
	sell = "ask"

	// 시장가 주문(매수)
	marketPriceBuy = "price"
	// 시장가 주문(매도)
	marketPriceSell = "market"

	// 주문하기 URL
	orderUrl = "https://api.upbit.com/v1/orders"
)

type Order struct {
	Market     string `json:"market" binding:"required" example:"KRW-BTC"` // 마켓 ID
	Side       string `json:"side" binding:"required"`                     // 주문 종류
	Volume     string `json:"volume"`                                      // 주문량 (지정가, 시장가 매도 시 필수)
	Price      string `json:"price"`
	OrdType    string `json:"ord_type"`
	Identifier string `json:"identifier"`
}

type Market string

// 설정한 액수만큼 코인 구매
func (m Market) BidMarketPrice(amount string) {

	request(
		Order{
			string(m),
			buy,
			"",
			amount,
			marketPriceBuy,
			"",
		},
	)

}

func request(order Order) {

	body := url.Values{}

	body.Set("market", order.Market)
	body.Set("side", order.Side)
	body.Set("volume", order.Volume)
	body.Set("price", order.Price)
	body.Set("ord_type", order.OrdType)

	fmt.Println(body)

	query := body.Encode()

	authorizationToken := middlewares.CreateTokenWithParams(query)

	client := &http.Client{}

	fmt.Println("---")
	fmt.Println(query)
	fmt.Println("---")

	b, _ := json.Marshal(order)

	req, err := http.NewRequest(http.MethodPost, orderUrl, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}
	req.Header.Add("Authorization", authorizationToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respBody))
}
