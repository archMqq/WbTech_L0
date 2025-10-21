package kconsumer

import (
	"L0/internal/database/models"
	"encoding/json"
	"fmt"
	"log"
)

func orderToJson(msg []byte) (*models.Order, error) {
	var order models.Order

	err := json.Unmarshal(msg, &order)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("JSON deserializing error: %s", msg)
	}

	return &order, nil
}
