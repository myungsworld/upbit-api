package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"upbit-api/internal/middlewares"
	"upbit-api/internal/models"
)

const (
	GetAccountEndPoint   = "https://api.upbit.com/v1/accounts"
	OrderEndPoint        = "https://api.upbit.com/v1/orders"
	OrderDeleteEndPoint  = "https://api.upbit.com/v1/order"
	GetCandleDayEndPoint = "https://api.upbit.com/v1/candles/days"
)

func Request(endPoint string, body interface{}) interface{} {

	var method string
	var token string

	requestBody := url.Values{}

	switch endPoint {
	case GetAccountEndPoint:
		method = http.MethodGet
		token = middlewares.CreateTokenWithNoParams()
	case OrderEndPoint, OrderDeleteEndPoint:
		method = http.MethodPost
		switch order := body.(type) {
		case models.GetOrder:
			method = http.MethodGet
			requestBody.Add("uuid", order.Uuid)
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
		case models.LimitOrder:
			requestBody.Set("volume", order.Volume)
			requestBody.Set("price", order.Price)
			requestBody.Set("market", order.Market)
			requestBody.Set("side", order.Side)
			requestBody.Set("ord_type", order.OrdType)
		case models.OrderList:
			method = http.MethodGet
			requestBody.Set("state", order.State)
			for _, state := range order.States {
				requestBody.Add("states", state)
			}
		case models.CancelOrder:
			method = http.MethodDelete
			requestBody.Set("uuid", order.Uuid)
		}

		query := requestBody.Encode()
		token = middlewares.CreateTokenWithParams(query)

	default:
		method = http.MethodGet
	}

	client := &http.Client{}

	var req *http.Request
	var err error

	if body == nil {
		req, err = http.NewRequest(method, endPoint, nil)
	} else {
		switch order := body.(type) {
		case models.BidOrder, models.AskOrder, models.LimitOrder, models.CancelOrder:
			b, _ := json.Marshal(&order)
			req, err = http.NewRequest(method, endPoint, bytes.NewBuffer(b))
		case models.GetOrder, models.OrderList:
			req, err = http.NewRequest(method, fmt.Sprintf("%s?%s", endPoint, requestBody.Encode()), nil)
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

		// TODO : 이게 무슨 케이스인지 알아야함
		if strings.Contains(err.Error(), "server sent GOAWAY") {
			fmt.Println("----")
			fmt.Println("method:", resp.Request.Method, " url:", resp.Request.URL)
			fmt.Println(err)
			fmt.Println(resp)
			fmt.Println("----")
		}

		panic(err)
	}

	var respCode interface{}

	switch resp.StatusCode {

	case 200, 201:
		switch {
		case endPoint == GetAccountEndPoint:
			respCode = &models.Accounts{}
		case endPoint == OrderEndPoint:
			switch resp.Request.Method {
			case "POST":
				respCode = &models.RespOrder{}
			case "GET":
				respCode = &[]models.RespOrder{}

			}
		case endPoint == OrderDeleteEndPoint:
			respCode = &models.RespOrder{}

		case strings.Contains(endPoint, GetCandleDayEndPoint):
			respCode = &models.ResponseDay{}
		default:
			fmt.Println(string(respBody))
		}

	case 400:
		switch {
		// 매수 체결 금액 부족
		case strings.Contains(string(respBody), "insufficient_funds_bid"):
			respCode = &models.Response400Error{}
		// 매도 체결 금액 부족
		case strings.Contains(string(respBody), "insufficient_funds_ask"):
			respCode = &models.Response400Error{}
		default:
			panic(string(respBody))
			return nil
		}

	case 401:
		log.Fatalf("err: %d , %v", resp.StatusCode, string(respBody))

	case 404:
		switch {
		// 주문 대기 취소 실패
		case strings.Contains(string(respBody), "order_not_found"):
			respCode = &models.Response404Error{}
		default:
			log.Fatalf("err: %d , %v", resp.StatusCode, string(respBody))
		}
	case 429:
		respCode = &models.Response429Error{}

	case 500:
		fmt.Println(string(respBody))
		fmt.Println(resp.Request)
		respCode = &models.Response404Error{}
	case 503:
		fmt.Println(string(respBody))
		respCode = &models.Response503Error{}
	default:
		log.Fatalf("err: %d , %v", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, &respCode); err != nil {
		panic(err)
	}

	return respCode
}
