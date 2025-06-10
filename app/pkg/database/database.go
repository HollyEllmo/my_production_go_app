package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// DB представляет подключение к базе данных
type DB struct {
	*sql.DB
}

// New создает новое подключение к базе данных
func New(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: db}, nil
}

// Close закрывает подключение к базе данных
func (db *DB) Close() error {
	return db.DB.Close()
}
