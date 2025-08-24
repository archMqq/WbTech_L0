package database

import (
	"L0/internal/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type OrderRepository struct {
	db *sql.DB
}

func Init(cfg *config.Config) *OrderRepository {
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Postgre.Host, cfg.Postgre.Port, cfg.Postgre.User,
		cfg.Postgre.Password, cfg.Postgre.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error opening db connection: %s", err)
	}

	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &OrderRepository{db: db}
}
