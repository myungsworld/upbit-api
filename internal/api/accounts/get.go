package accounts

import (
	"encoding/json"
	"io"
	"net/http"
	"upbit-api/internal/middlewares"
	"upbit-api/internal/models"
)

// https://docs.upbit.com/reference/%EC%A0%84%EC%B2%B4-%EA%B3%84%EC%A2%8C-%EC%A1%B0%ED%9A%8C

// Index 전체 계좌 조회
func Index() models.Accounts {

	client := &http.Client{}

	token := middlewares.CreateTokenWithNoParams()

	req, err := http.NewRequest(http.MethodGet, "https://api.upbit.com/v1/accounts", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	accounts := models.Accounts{}

	if err = json.Unmarshal(body, &accounts); err != nil {
		panic(err)
	}

	return accounts
}
