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
	marketLimitBuy = "limit"

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
func (m Market) AskMarketPrice(volume string) {
	api.Request("https://api.upbit.com/v1/orders", models.AskOrder{
		Market:  string(m),
		OrdType: marketPriceSell,
		Volume:  volume,
		Side:    sell,
	})
}

// BidMarketLimit 지정가 매수 주문
func (m Market) BidMarketLimit(bidPrice, bidVolume string) bool {

	if resp := api.Request("https://api.upbit.com/v1/orders", models.LimitOrder{
		Market:  string(m),
		OrdType: marketLimitBuy,
		Price:   bidPrice,
		Side:    buy,
		Volume:  bidVolume,
	}); resp != nil {

		switch resp.(type) {
		case *models.RespOrder:
			fmt.Println(string(m), "매수대기 체결", bidPrice)
			return true
		case *models.Response400Error:
			fmt.Println(string(m), "주문금액 부족으로 인해 해당코인 비활성화")
			return true
		}

	}

	return false

}
