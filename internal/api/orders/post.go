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

	// 주문하기 URL
	orderUrl = "https://api.upbit.com/v1/orders"
)

type Market string

// 설정한 액수만큼 코인 구매
func (m Market) BidMarketPrice(amount string) {

	api.Request("https://api.upbit.com/v1/orders", models.BidOrder{
		Market:  string(m),
		OrdType: marketPriceBuy,
		Price:   amount,
		Side:    buy,
	})

	//request(
	//	BidOrder{
	//		string(m),
	//		marketPriceBuy,
	//		amount,
	//		buy,
	//	},
	//)

}

// 시장가 매도
func (m Market) AskMarketPrice(volume string) {
	api.Request("https://api.upbit.com/v1/orders", models.AskOrder{
		Market:  string(m),
		OrdType: marketPriceSell,
		Volume:  volume,
		Side:    sell,
	})
}
