package orders

import (
	"fmt"
	"upbit-api/internal/api"
	"upbit-api/internal/models"
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
	// 지정가 주문
	marketLimit = "limit"

	// 주문하기 URL
	orderUrl = "https://api.upbit.com/v1/orders"
)

type Market string

// BidMarketPrice 시장가 매수
func (m Market) BidMarketPrice(amount string) {

	api.Request("https://api.upbit.com/v1/orders", models.BidOrder{
		Market:  string(m),
		OrdType: marketPriceBuy,
		Price:   amount,
		Side:    buy,
	})

}

// AskMarketPrice 시장가 매도
func (m Market) AskMarketPrice(volume string) string {
	if resp := api.Request("https://api.upbit.com/v1/orders", models.AskOrder{
		Market:  string(m),
		OrdType: marketPriceSell,
		Volume:  volume,
		Side:    sell,
	}); resp != nil {
		switch resp.(type) {
		case *models.RespOrder:
			uuid := resp.(*models.RespOrder).Uuid
			return uuid
		default:
			return "시장가 매도"
		}
	}
	return "시장가 매도"
}

// BidMarketLimit 지정가 매수 주문
func (m Market) BidMarketLimit(bidPrice, bidVolume string) string {

	if resp := api.Request("https://api.upbit.com/v1/orders", models.LimitOrder{
		Market:  string(m),
		OrdType: marketLimit,
		Price:   bidPrice,
		Side:    buy,
		Volume:  bidVolume,
	}); resp != nil {

		switch resp.(type) {
		case *models.RespOrder:
			uuid := resp.(*models.RespOrder).Uuid
			fmt.Println(string(m), "매수대기 체결", bidPrice)
			return uuid
		case *models.Response400Error:
			fmt.Println(string(m), "주문금액 부족으로 인해 해당코인 비활성화")
			return "0"
		case *models.Response429Error:
			res := resp.(*models.Response429Error)
			fmt.Println(res)
			panic(res)
		default:
			panic(resp)
			return "-1"
		}

	}

	return "-1"
}

// 지정가 매도 주문
func (m Market) AskMarketLimit(askPrice, askVolume string) string {

	if resp := api.Request("https://api.upbit.com/v1/orders", models.LimitOrder{
		Market:  string(m),
		Side:    sell,
		Volume:  askVolume,
		Price:   askPrice,
		OrdType: marketLimit,
	}); resp != nil {
		switch resp.(type) {
		case *models.RespOrder:
			uuid := resp.(*models.RespOrder).Uuid
			fmt.Println(string(m), "매도대기 체결", askPrice)
			return uuid
		case *models.Response400Error:
			fmt.Println(string(m), "매도대기 체결 금액 부족")
			panic("매도 대기 체결 금액 부족")
			return "0"
		default:
			panic(resp)
			return "-1"
		}
	}

	return "-1"
}
