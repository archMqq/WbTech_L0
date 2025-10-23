package fake

import (
	"L0/internal/config"
	"L0/internal/database/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/segmentio/kafka-go"
)

func InitWriter(cfg *config.Config) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{cfg.Kafka.Broker},
		Topic:   cfg.Kafka.Topic,
	})
}

func StartFaking(ctx context.Context, k *kafka.Writer) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fakeObj := controlledFakeOrder()

			jsonData, err := json.Marshal(fakeObj)
			if err != nil {
				log.Printf("Error marshaling order: %v", err)
				continue
			}

			message := kafka.Message{
				Key:   []byte(fakeObj.OrderUID),
				Value: []byte(jsonData),
			}

			err = k.WriteMessages(ctx, message)
			if err != nil {
				log.Printf("error writing fake to kafka: %s", err)
			} else {
				log.Printf("wrote message with uid: %s: ", fakeObj.OrderUID)
			}
		}
	}
}

func controlledFakeOrder() *models.Order {
	orderUID := "test_" + gofakeit.DigitN(10)

	return &models.Order{
		OrderUID:          orderUID,
		TrackNumber:       "TRACK_" + gofakeit.DigitN(5),
		Entry:             "ENTRY",
		Delivery:          controlledFakeDelivery(orderUID),
		Payment:           controlledFakePayment(orderUID),
		Items:             []models.Item{controlledFakeItem(orderUID)},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "customer_" + gofakeit.DigitN(5),
		DeliveryService:   "test_delivery",
		Shardkey:          "1",
		SmID:              1,
		DateCreated:       time.Now().UTC(),
		OofShard:          "1",
	}
}

func controlledFakeDelivery(orderUID string) models.Delivery {
	return models.Delivery{
		OrderUID: orderUID,
		Name:     "Test User",
		Phone:    "+1234567890",
		Zip:      "12345",
		City:     "Test City",
		Address:  "123 Test Street",
		Region:   "Test Region",
		Email:    "test@example.com",
	}
}

func controlledFakePayment(orderUID string) models.Payment {
	return models.Payment{
		Transaction:  orderUID,
		RequestID:    "",
		Currency:     "USD",
		Provider:     "test_provider",
		Amount:       1000,
		PaymentDt:    time.Now().Unix(),
		Bank:         "test_bank",
		DeliveryCost: 500,
		GoodsTotal:   500,
		CustomFee:    0,
	}
}

func controlledFakeItem(orderUID string) models.Item {
	return models.Item{
		OrderUID:    orderUID,
		ChrtID:      1234567,
		TrackNumber: "TRACK_12345",
		Price:       100,
		Rid:         "rid_123456789",
		Name:        "Test Product",
		Sale:        10,
		Size:        "M",
		TotalPrice:  90,
		NmID:        7654321,
		Brand:       "Test Brand",
		Status:      200,
	}
}
