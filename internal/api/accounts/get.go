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

	//client := &http.Client{}
	//
	//token := middlewares.CreateTokenWithNoParams()
	//
	//req, err := http.NewRequest(http.MethodGet, "https://api.upbit.com/v1/accounts", nil)
	//if err != nil {
	//	panic(err)
	//}
	//
	//req.Header.Add("Authorization", token)
	//req.Header.Add("Content-Type", "application/json")
	//
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//
	//accounts := models.Accounts{}
	//
	//if err = json.Unmarshal(body, &accounts); err != nil {
	//	panic(err)
	//}
	//
	//return accounts
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
