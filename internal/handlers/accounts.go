package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
	"upbit-api/internal/middlewares"
)

func Accounts(c *gin.Context) {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	token := middlewares.Authentication()

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

	fmt.Println(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

}
