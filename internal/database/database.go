package database

import (
	"L0/internal/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func Init(cfg *config.Config) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		cfg.Postgre.Host, cfg.Postgre.Port, cfg.Postgre.User,
		cfg.Postgre.Password, cfg.Postgre.DBName, cfg.Postgre.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error opening db connection: %s", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}
