package orders

import (
	"fmt"
	"upbit-api/internal/api"
	"upbit-api/internal/models"
)

func Cancel(uuid string) {

	result := api.Request(api.OrderDeleteEndPoint, models.CancelOrder{
		Uuid: uuid,
	})

	switch result.(type) {
	case *models.RespOrder:
		respOrder := result.(*models.RespOrder)
		fmt.Println(respOrder.Market, "대기 주문 취소 완료")
	case *models.Response404Error:
		//respError := result.(*models.Response404Error)
		fmt.Println("취소할 대기주문이 없음")
	default:
		panic("다른 케이스")
	}

}
