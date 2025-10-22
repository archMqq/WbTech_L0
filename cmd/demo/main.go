package main

import (
	"L0/internal/config"
	"L0/internal/database"
	"L0/internal/handler"
	kconsumer "L0/internal/kafka"
	"L0/internal/repository"
	"L0/internal/services"
	"L0/internal/validation"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	db := database.Init(cfg)
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)
	orderCache := repository.NewInMemoryCache(time.Second*5, time.Second*10)
	validator := validation.NewValidator()

	service := services.NewOrderService(*orderRepo, *orderCache, validator)

	reader := kconsumer.InitReader(cfg)
	go kconsumer.Start(reader, service)

	orderHandler := handler.NewOrderHandler(service)

	router := mux.NewRouter()
	router.Use(corsMiddleware)
	router.HandleFunc("/order", orderHandler.GetOrder).Methods("GET")

	log.Printf("Server starting on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
