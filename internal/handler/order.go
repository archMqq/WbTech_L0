package handler

import (
	"L0/internal/services"
	"encoding/json"
	"net/http"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (rh *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	order, err := rh.service.GetOrder(id)
	if err != nil {
		http.Error(w, "Object was not found", http.StatusNotFound)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(order)
}
