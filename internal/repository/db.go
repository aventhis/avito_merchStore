// Package repository содержит функции для подключения к базе данных.
package repository

import (
	"avito_merchStore/internal/config"
	"database/sql"
	"fmt"

	// Импорт драйвера PostgreSQL для регистрации его в пакете database/sql.
	_ "github.com/lib/pq"
)

// NewPostgresDB устанавливает подключение к базе данных PostgreSQL,
// используя параметры из конфигурации. Возвращает подключение или ошибку.
func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
