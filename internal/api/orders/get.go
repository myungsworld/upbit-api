package orders

import (
	"fmt"
	"upbit-api/internal/api"
	"upbit-api/internal/models"
)

// 체결 대기 리스트
func WaitList() *[]models.RespOrder {
	result := api.Request(api.OrderEndPoint, models.OrderList{
		State: "wait",
	})

	switch result.(type) {
	case *[]models.RespOrder:
		return result.(*[]models.RespOrder)
	}

	return nil
}

// 체결 완료 리스트
func GetTodayDoneList() *[]models.RespOrder {
	result := api.Request(api.OrderEndPoint, models.OrderList{
		Market: "KRW-IQ",
		State:  "done",
		States: []string{"done"},
	})

	switch result.(type) {
	case *[]models.RespOrder:
		return result.(*[]models.RespOrder)
	default:
		panic("이상함")
	}

	return nil
}

func Get(uuid string) {
	result := api.Request(api.OrderDeleteEndPoint, models.GetOrder{
		Uuid: uuid,
	})

	switch result.(type) {
	case *models.RespOrder:
		fmt.Println(result.(*models.RespOrder))
	default:
		panic("아님")
	}
}
