package models

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required,min=10,max=20"`
	TrackNumber       string    `json:"track_number" validate:"required,min=5,max=15"`
	Entry             string    `json:"entry" validate:"required,min=1,max=10"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" validate:"required,min=2,max=6"`
	InternalSignature string    `json:"internal_signature" validate:"max=50"`
	CustomerID        string    `json:"customer_id" validate:"required,min=1,max=25"`
	DeliveryService   string    `json:"delivery_service" validate:"required,min=1,max=30"`
	Shardkey          string    `json:"shardkey" validate:"required,min=1,max=5"`
	SmID              int       `json:"sm_id" validate:"required,min=0"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required,min=1,max=5"`
}

type Delivery struct {
	OrderUID string `json:"-"`
	Name     string `json:"name" validate:"required,min=1,max=50"`
	Phone    string `json:"phone" validate:"required,min=5,max=50"`
	Zip      string `json:"zip" validate:"required,min=1,max=50"`
	City     string `json:"city" validate:"required,min=1,max=50"`
	Address  string `json:"address" validate:"required,min=1"`
	Region   string `json:"region" validate:"required,min=1,max=50"`
	Email    string `json:"email" validate:"required,email,max=50"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required,min=10,max=20"`
	RequestID    string `json:"request_id" validate:"max=20"`
	Currency     string `json:"currency" validate:"required,min=1,max=10"`
	Provider     string `json:"provider" validate:"required,min=1,max=100"`
	Amount       int    `json:"amount" validate:"required,min=0"`
	PaymentDt    int64  `json:"payment_dt" validate:"required,min=0"`
	Bank         string `json:"bank" validate:"required,min=1,max=100"`
	DeliveryCost int    `json:"delivery_cost" validate:"min=0"`
	GoodsTotal   int    `json:"goods_total" validate:"required,min=0"`
	CustomFee    int    `json:"custom_fee" validate:"min=0"`
}

type Item struct {
	OrderUID    string `json:"-"`
	ChrtID      int    `json:"chrt_id" validate:"required,min=1"`
	TrackNumber string `json:"track_number" validate:"required,min=5,max=15"`
	Price       int    `json:"price" validate:"required,min=0"`
	Rid         string `json:"rid" validate:"required,min=1,max=25"`
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Sale        int    `json:"sale" validate:"min=0"`
	Size        string `json:"size" validate:"required,min=1,max=50"`
	TotalPrice  int    `json:"total_price" validate:"required,min=0"`
	NmID        int    `json:"nm_id" validate:"required,min=1"`
	Brand       string `json:"brand" validate:"required,min=1,max=50"`
	Status      int    `json:"status" validate:"required"`
}
