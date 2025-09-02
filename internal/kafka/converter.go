package kconsumer

import (
	"L0/internal/database/models"
	"encoding/json"
)

func orderToJson(msg []byte) *models.Order {
	var order *models.Order

	err := json.Unmarshal(msg, order)
	if err != nil {
		return nil
	}

	return order
}
