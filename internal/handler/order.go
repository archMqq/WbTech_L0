package handler

import (
	"L0/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	repo *repository.OrderRepository
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

func (rh *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	order := rh.repo.GetCached(id)
	if order == nil {
		var err error

		order, err = rh.repo.GetOrderByID(id)
		if err != nil {
			http.Error(w, "Object was not found", http.StatusNotFound)
		}

		go rh.repo.SaveCache(order, id)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(order)
}
