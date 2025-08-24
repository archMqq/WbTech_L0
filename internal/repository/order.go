package repository

import (
	"L0/internal/database/models"
	"database/sql"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetOrderByID(orderUID string) (*models.Order, error) {
	order, err := r.getOrder(orderUID)
	if err != nil {
		return nil, err
	}

	delivery, err := r.getDelivery(orderUID)
	if err != nil {
		return nil, err
	}
	order.Delivery = *delivery

	payment, err := r.getPayment(order.TrackNumber)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment

	items, err := r.getItems(orderUID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (r *OrderRepository) getOrder(orderUID string) (*models.Order, error) {
	query := `
        SELECT order_uid, track_number, entry, locale, internal_signature, 
               customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders 
        WHERE order_uid = $1
	`

	var order models.Order
	err := r.db.QueryRow(query, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) getDelivery(orderUID string) (*models.Delivery, error) {
	query := `
		SELECT name, phone, zip, city, address, region, email
        FROM delivery 
        WHERE order_uid = $1
	`

	var delivery models.Delivery
	delivery.OrderUID = orderUID

	err := r.db.QueryRow(query, orderUID).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)

	if err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *OrderRepository) getPayment(transaction string) (*models.Payment, error) {
	query := `
		SELECT transaction, request_id, currency, provider, amount, 
			   payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments
		WHERE transaction = $1
	`

	var payment models.Payment
	err := r.db.QueryRow(query, transaction).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *OrderRepository) getItems(orderUID string) ([]models.Item, error) {
	query := `
        SELECT chrt_id, track_number, price, rid, name, sale, size, 
               total_price, nm_id, brand, status
        FROM items 
        WHERE order_uid = $1
    `

	rows, err := r.db.Query(query, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		item.OrderUID = orderUID

		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
