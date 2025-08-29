package config

import (
	"database/sql"
	"fmt"
)

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

func Connect() (*sql.DB, error) {
	dsn := "postgres://user:pass@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %w", err)
	}
	return db, nil
}
