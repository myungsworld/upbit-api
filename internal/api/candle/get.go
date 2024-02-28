package candle

import (
	"fmt"
	"upbit-api/internal/api"
)

type Market string

func (m Market) Min() {

	api.Request(fmt.Sprintf("https://api.upbit.com/v1/candles/minutes/1?market=%v&count=1", string(m)), nil)

}

func (m Market) Day(count int) {
	api.Request(fmt.Sprintf("https://api.upbit.com/v1/candles/days?market=%v&count=%d", string(m), count), nil)
}
