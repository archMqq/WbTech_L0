package repository

import "L0/internal/database/models"

type OrderCache interface {
	Set(id string, order *models.Order)
	Get(id string) *models.Order
}

type OrderRepo interface {
	GetOrderByID(orderUID string) (*models.Order, error)
	SaveOrder(order *models.Order) error
}
