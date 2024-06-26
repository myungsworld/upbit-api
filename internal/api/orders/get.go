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
func GetDoneList() *[]models.RespOrder {
	result := api.Request(api.OrderEndPoint, models.OrderList{
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

func Get(uuid string) *models.RespOrder {
	result := api.Request(api.OrderDeleteEndPoint, models.GetOrder{
		Uuid: uuid,
	})

	switch result.(type) {
	case *models.RespOrder:
		return result.(*models.RespOrder)
	default:
		fmt.Println("------")
		fmt.Println("uuid 조회 정상적이지 않은 값")
		fmt.Println(result)
		fmt.Println("------")
		return nil
	}

}
