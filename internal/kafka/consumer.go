package kconsumer

import (
	"L0/internal/config"
	"L0/internal/services"
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func InitReader(cfg *config.Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Kafka.Broker},
		Topic:   cfg.Kafka.Topic,
	})
}

func Start(reader *kafka.Reader, service *services.OrderService) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading kafka message: %s", err)
			break
		}

		order, err := orderToJson(m.Value)
		if err != nil {
			log.Printf("unknown message format: %s", err)
		} else if err := service.SaveOrder(order); err != nil {
			log.Fatalf("error request to db: %s", err)
		}
	}

	if err := reader.Close(); err != nil {
		log.Printf("error closing consumer")
	}
}
