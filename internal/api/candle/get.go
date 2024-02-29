package candle

import (
	"fmt"
	"upbit-api/internal/api"
	"upbit-api/internal/models"
)

type Market string

func (m Market) Min() {

	api.Request(fmt.Sprintf("https://api.upbit.com/v1/candles/minutes/1?market=%v&count=1", string(m)), nil)

}

func (m Market) Day(count int) *models.ResponseDay {
	result := api.Request(fmt.Sprintf("https://api.upbit.com/v1/candles/days?market=%v&count=%d", string(m), count), nil)
	responseDays := result.(*models.ResponseDay)
	return responseDays
}
