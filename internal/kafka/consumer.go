package kconsumer

import (
	"L0/internal/config"
	"L0/internal/repository"
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

func Start(reader *kafka.Reader, repo *repository.OrderRepository) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading kafka message: %s", err)
			break
		}

		order := orderToJson(m.Value)
		if order == nil {
			log.Print("error json convert")
		}

		if err := repo.SaveOrder(order); err != nil {
			log.Fatalf("error request to db: %s", err)
		}
	}

	if err := reader.Close(); err != nil {
		log.Printf("error closing consumer")
	}
}
