package repository

import "L0/internal/database/models"

type OrderCache interface {
	Set(id string, order *models.Order)
	Get(id string) *models.Order
}

type MemoryCache interface {
	StartCollector()
}

type OrderRepo interface {
	GetAllOrders() (map[string]*models.Order, error)
	GetOrderByID(orderUID string) (*models.Order, error)
	SaveOrder(order *models.Order) error
	OrderExists(orderUID string) (bool, error)
}
