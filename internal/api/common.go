package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"upbit-api/internal/middlewares"
	"upbit-api/internal/models"
)

const (
	GetAccountEndPoint = "https://api.upbit.com/v1/accounts"
	OrderEndPoint      = "https://api.upbit.com/v1/orders"
)

func Request(endPoint string, body interface{}) interface{} {

	var method string
	var token string

	requestBody := url.Values{}

	switch endPoint {
	case GetAccountEndPoint:
		method = http.MethodGet
		token = middlewares.CreateTokenWithNoParams()
	case OrderEndPoint:
		method = http.MethodPost
		switch order := body.(type) {
		case models.BidOrder:
			requestBody.Set("price", order.Price)
			requestBody.Set("market", order.Market)
			requestBody.Set("side", order.Side)
			requestBody.Set("ord_type", order.OrdType)
		case models.AskOrder:
			requestBody.Set("volume", order.Volume)
			requestBody.Set("market", order.Market)
			requestBody.Set("side", order.Side)
			requestBody.Set("ord_type", order.OrdType)
		}

		query := requestBody.Encode()
		token = middlewares.CreateTokenWithParams(query)
	}

	client := &http.Client{}

	var req *http.Request
	var err error

	if body == nil {
		req, err = http.NewRequest(method, endPoint, nil)
	} else {
		switch order := body.(type) {
		case models.BidOrder, models.AskOrder:
			fmt.Println(order)
			b, _ := json.Marshal(&order)
			req, err = http.NewRequest(method, endPoint, bytes.NewBuffer(b))
		}
	}

	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result := respHandler(endPoint, resp)

	return result

}

func respHandler(endPoint string, resp *http.Response) interface{} {

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)

	var respCode interface{}

	switch resp.StatusCode {
	case 401:
		log.Fatalf("err: %d , %v", resp.StatusCode, string(respBody))
	case 200, 201:
		switch endPoint {
		case GetAccountEndPoint:
			respCode = &models.Accounts{}
		case OrderEndPoint:
			respCode = &models.RespOrder{}
		}
	default:
		log.Fatalf("err: %d , %v", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, &respCode); err != nil {
		panic(err)
	}

	return respCode
}