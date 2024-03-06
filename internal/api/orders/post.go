package orders

import (
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
func (m Market) BidMarketLimit(amount string) {

	var price string
	var volume string

	// 현금 = 주문가격 * 주문량
	// amount = price * volume

	api.Request("https://api.upbit.com/v1/orders", models.LimitOrder{
		Market:  string(m),
		OrdType: marketPriceBuy,
		Price:   price,
		Side:    buy,
		Volume:  volume,
	})

}
