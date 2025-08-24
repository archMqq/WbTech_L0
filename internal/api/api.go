package api

import (
	"L0/internal/api/handlers/order"
	"L0/internal/config"
	"L0/internal/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var rep *database.OrderRepository

func Init(cfg *config.Config, repository *database.OrderRepository) {
	rep = repository

	r := mux.NewRouter()
	r.HandleFunc("/order", order.GetOrder).Methods("GET")

	log.Fatal(http.ListenAndServe(cfg.Port, r))
}
