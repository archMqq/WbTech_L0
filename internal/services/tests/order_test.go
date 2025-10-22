package services__test

import (
	"L0/internal/database/models"
	"L0/internal/repository/mocks"
	"L0/internal/services"
	"L0/internal/validation"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createValidTestOrder() *models.Order {
	return &models.Order{
		OrderUID:          "test1234567890",
		TrackNumber:       "TRACK1234567",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test_customer",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now().Add(time.Second * -1).UTC(),
		OofShard:          "1",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "test1234567890",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "TRACK1234567",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
	}
}

func TestOrderService_CreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	mockCache := mocks.NewMockOrderCache(ctrl)
	validator := validation.NewValidator()

	service := services.NewOrderService(mockRepo, mockCache, validator)

	order := createValidTestOrder()

	mockRepo.EXPECT().
		OrderExists("test1234567890").
		Return(false, nil)

	mockRepo.EXPECT().
		GetAllOrders().
		Return(map[string]*models.Order{"test1234567890": order}, nil)

	mockRepo.EXPECT().
		SaveOrder(gomock.Any()).
		Return(nil)

	mockCache.EXPECT().
		Set("test1234567890", order).
		AnyTimes()

	err := service.SaveOrder(order)

	assert.NoError(t, err)
}

func TestOrderService_CreateOrder_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	mockCache := mocks.NewMockOrderCache(ctrl)
	validator := validation.NewValidator()

	service := services.NewOrderService(mockRepo, mockCache, validator)

	order := createValidTestOrder()

	mockRepo.EXPECT().
		GetAllOrders().
		Return(map[string]*models.Order{}, nil)

	mockRepo.EXPECT().
		OrderExists("test1234567890").
		Return(true, nil)
	err := service.SaveOrder(order)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestOrderService_GetOrder_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	mockCache := mocks.NewMockOrderCache(ctrl)
	validator := validation.NewValidator()

	service := services.NewOrderService(mockRepo, mockCache, validator)

	expectedOrder := createValidTestOrder()

	mockCache.EXPECT().
		Get("test123").
		Return(expectedOrder)

	order, err := service.GetOrder("test123")

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestOrderService_GetOrder_CacheMiss(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	mockCache := mocks.NewMockOrderCache(ctrl)
	validator := validation.NewValidator()

	service := services.NewOrderService(mockRepo, mockCache, validator)

	expectedOrder := createValidTestOrder()

	mockCache.EXPECT().
		Get("test123").
		Return(nil)

	mockRepo.EXPECT().
		GetOrderByID("test123").
		Return(expectedOrder, nil)

	mockRepo.EXPECT().
		GetAllOrders().
		Return(map[string]*models.Order{}, nil).
		AnyTimes()

	mockCache.EXPECT().
		Set("test123", expectedOrder)

	order, err := service.GetOrder("test123")

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}
