package orders

import (
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
		State: "cancel",
		//States: []string{"cancel"},
	})

	switch result.(type) {
	case *[]models.RespOrder:
		return result.(*[]models.RespOrder)
	default:
		panic("이상함")
	}

	return nil
}
