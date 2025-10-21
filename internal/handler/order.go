package handler

import (
	"L0/internal/repository"
	"encoding/json"
	"net/http"
)

type OrderHandler struct {
	repo *repository.OrderRepository
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

func (rh *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	repo := rh.repo

	order := repo.Cache.Get(id)
	if order == nil {
		var err error

		order, err = repo.GetOrderByID(id)
		if err != nil {
			http.Error(w, "Object was not found", http.StatusNotFound)
		}

		go repo.Cache.Set(order, id)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(order)
}
