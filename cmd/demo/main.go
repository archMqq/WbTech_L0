package main

import (
	"L0/internal/config"
	"L0/internal/database"
	"L0/internal/handler"
	kconsumer "L0/internal/kafka"
	"L0/internal/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	db := database.Init(cfg)
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)

	reader := kconsumer.InitReader(cfg)
	go kconsumer.Start(reader, orderRepo)

	orderHandler := handler.NewOrderHandler(orderRepo)

	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")

	log.Printf("Server starting on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, router))
}
