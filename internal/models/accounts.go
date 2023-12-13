package models

type Accounts []struct {
	Currency     string `json:"currency"`      // 가상화폐
	Balance      string `json:"balance"`       // 주문가능 금액/수량
	Locked       string `json:"locked"`        // 주문중 묶여있는 금액/수량
	AvgBuyPrice  string `json:"avg_buy_price"` // 매수 평균가
	UnitCurrency string `json:"unit_currency"` // 평단가 기준화폐
}
