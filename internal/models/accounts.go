package models

type Accounts []struct {
	Currency    string `json:"currency"`
	Balance     string `json:"balance"`
	Locked      string `json:"locked"`
	AvgBuyPrice string `json:"avg_buy_price"`
}

func (a Accounts) GetAvailableKRW() string {
	for _, account := range a {
		if account.Currency == "KRW" {
			return account.Balance
		}
	}
	return ""
}
