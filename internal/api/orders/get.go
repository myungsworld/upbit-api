package orders

import (
	"upbit-api/internal/api"
	"upbit-api/internal/models"
)

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
