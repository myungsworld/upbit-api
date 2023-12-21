package accounts

import (
	"strconv"
	"upbit-api/internal/api"
	"upbit-api/internal/models"
)

// https://docs.upbit.com/reference/%EC%A0%84%EC%B2%B4-%EA%B3%84%EC%A2%8C-%EC%A1%B0%ED%9A%8C

// Index 전체 계좌 조회
func index() *models.Accounts {

	result := api.Request("https://api.upbit.com/v1/accounts", nil)
	accounts := result.(*models.Accounts)

	return accounts
}

// 현재 사용가능한 한화 금액
func GetAvailableKRW() int {

	accounts := index()

	for _, account := range *accounts {
		if account.Currency == "KRW" {
			floatValue, _ := strconv.ParseFloat(account.Balance, 64)

			return int(floatValue)
		}
	}
	return 0
}

func Get() models.Accounts {
	accounts := index()

	return *accounts
}
