package middlewares

type Order struct {
	Market     string `json:"market" binding:"required" example:"KRW-BTC"`
	Side       string `json:"side" binding:"required"`
	Volume     string `json:"volume"`
	Price      string `json:"price"`
	OrdType    string `json:"ord_type"`
	Identifier string `json:"identifier"`
}
