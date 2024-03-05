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

	var day *models.ResponseDay
	//var err503 *models.Response503Error
	//var err429 *models.Response429Error

	switch result.(type) {
	case *models.Response503Error:
		//err503 = result.(*models.Response503Error)
		//fmt.Println(err503)
		return nil
	case *models.Response429Error:
		//err429 = result.(*models.Response429Error)
		//fmt.Println(err429)
		return nil
	default:
		day = result.(*models.ResponseDay)
		//fmt.Println(day)
		return day
	}

}
