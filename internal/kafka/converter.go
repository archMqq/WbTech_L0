package kconsumer

import (
	"L0/internal/database/models"
	"encoding/json"
	"log"
)

func orderToJson(msg []byte) *models.Order {
	var order models.Order

	err := json.Unmarshal(msg, &order)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &order
}
