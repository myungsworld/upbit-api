package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func AllMarketCodes(c *gin.Context) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.upbit.com/v1/market/all?isDetails=true", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	c.String(resp.StatusCode, string(body))

}
