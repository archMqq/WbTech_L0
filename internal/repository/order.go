package repository

import (
	"L0/internal/database/models"
	"database/sql"
	"fmt"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	repo := &OrderRepository{db: db}

	return repo
}

func (r OrderRepository) GetAllOrders() (map[string]*models.Order, error) {
	orderUIDs, err := r.getAllOrderUIDs()
	if err != nil {
		return nil, err
	}

	orders := make(map[string]*models.Order)
	for _, orderUID := range orderUIDs {
		order, err := r.GetOrderByID(orderUID)
		if err != nil {
			return nil, err
		}
		orders[orderUID] = order
	}

	return orders, nil
}

func (r OrderRepository) getAllOrderUIDs() ([]string, error) {
	query := "SELECT order_uid FROM orders ORDER BY date_created DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderUIDs []string
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return nil, err
		}
		orderUIDs = append(orderUIDs, orderUID)
	}

	return orderUIDs, nil
}

func (r OrderRepository) GetOrderByID(orderUID string) (*models.Order, error) {
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

func (r OrderRepository) getOrder(orderUID string) (*models.Order, error) {
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

func (r OrderRepository) getDelivery(orderUID string) (*models.Delivery, error) {
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

func (r OrderRepository) getPayment(transaction string) (*models.Payment, error) {
	query := `
		SELECT transaction, request_id, currency, provider, amount, 
			   payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
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

func (r OrderRepository) getItems(orderUID string) ([]models.Item, error) {
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

func (r OrderRepository) SaveOrder(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := r.saveOrderTx(tx, order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	if err := r.saveDeliveryTx(tx, &order.Delivery, order.OrderUID); err != nil {
		return fmt.Errorf("failed to save delivery: %w", err)
	}

	if err := r.savePaymentTx(tx, &order.Payment); err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}

	if err := r.saveItemsTx(tx, order.Items, order.OrderUID); err != nil {
		return fmt.Errorf("failed to save items: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r OrderRepository) saveOrderTx(tx *sql.Tx, order *models.Order) error {
	query := `
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
                          customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

	_, err := tx.Exec(query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)

	return err
}

func (r OrderRepository) saveDeliveryTx(tx *sql.Tx, delivery *models.Delivery, orderId string) error {
	query := `
        INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	_, err := tx.Exec(query,
		orderId,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)

	return err
}

func (r OrderRepository) savePaymentTx(tx *sql.Tx, payment *models.Payment) error {
	query := `
        INSERT INTO payment (transaction, request_id, currency, provider, amount, 
                           payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	_, err := tx.Exec(query,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)

	return err
}

func (r OrderRepository) saveItemsTx(tx *sql.Tx, items []models.Item, orderId string) error {
	query := `
        INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, 
                         sale, size, total_price, nm_id, brand, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `

	for _, item := range items {
		_, err := tx.Exec(query,
			orderId,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item %d: %w", item.ChrtID, err)
		}
	}

	return nil
}

func (r OrderRepository) OrderExists(orderUID string) (bool, error) {
	query := `SELECT * FROM orders WHERE order_uid = $1`

	rows, err := r.db.Query(query, orderUID)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}
