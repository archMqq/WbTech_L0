package services

import (
	"L0/internal/database/models"
	"L0/internal/repository"
	"L0/internal/validation"
	"fmt"
	"log"
)

type OrderService struct {
	repo      repository.OrderRepo
	cache     repository.OrderCache
	validator *validation.OrderValidator
}

func NewOrderService(
	repo repository.OrderRepo,
	cache repository.OrderCache,
	validator *validation.OrderValidator,
) *OrderService {
	if cache != nil {
		if obj, ok := cache.(repository.MemoryCache); ok {
			obj.StartCollector()
		}
	}

	service := &OrderService{
		repo:      repo,
		cache:     cache,
		validator: validator,
	}

	go service.loadCache()

	return service
}

func (s *OrderService) loadCache() {
	orders, err := s.repo.GetAllOrders()
	if err != nil {
		log.Printf("error cache loading: %s", err)
		return
	}

	for k, v := range orders {
		s.cache.Set(k, v)
	}
}

func (s *OrderService) SaveOrder(order *models.Order) error {
	if err := s.validator.Validate(order); err != nil {
		return fmt.Errorf("validation error: %s", err)
	}

	exists, err := s.repo.OrderExists(order.OrderUID)
	if err != nil {
		return fmt.Errorf("order checking failed: %s", err)
	}
	if exists {
		return fmt.Errorf("order with UID already exists: %s", order.OrderUID)
	}

	if err := s.repo.SaveOrder(order); err != nil {
		return fmt.Errorf("order saving error: %s", err)
	}

	s.cache.Set(order.OrderUID, order)

	return nil
}

func (s *OrderService) GetOrder(orderUID string) (*models.Order, error) {
	order := s.cache.Get(orderUID)
	if order != nil {
		return order, nil
	}

	order, err := s.repo.GetOrderByID(orderUID)
	if err != nil {
		return nil, fmt.Errorf("getting order error: %s", err)
	}

	s.cache.Set(orderUID, order)
	return order, nil
}
