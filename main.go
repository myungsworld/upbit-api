package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"upbit-api/config"
	"upbit-api/internal/handlers"
)

func main() {

	config.Init()

	fmt.Println(len(config.Coins))

	r := gin.Default()
	rAPI := r.Group("/api")

	rAPI.GET("/accounts", handlers.Accounts)
	rAPI.GET("/마켓코드조회", handlers.AllMarketCodes)
	r.Run()

}
