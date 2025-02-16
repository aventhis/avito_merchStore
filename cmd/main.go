// Package main является точкой входа в приложение Avito Merch Store.
package main

import (
	"avito_merchStore/internal/config"
	"avito_merchStore/internal/repository"
	"avito_merchStore/internal/routes"
	"avito_merchStore/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

// main загружает конфигурацию, инициализирует подключение к базе данных, сервисы и роутер,
// а затем запускает HTTP-сервер.
func main() {
	// Загружаем конфигурацию из переменных окружения
	cfg := config.LoadConfig()

	// Инициализируем подключение к БД PostgreSQL
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Инициализируем сервисы
	authService := service.NewAuthService(db, cfg.JWTSecret)
	merchService := service.NewMerchService(db)
	coinService := service.NewCoinService(db)

	// Создаем роутер с использованием Gin
	router := gin.Default()
	// Регистрируем маршруты, передаем также JWT-секрет из конфига
	routes.RegisterRoutes(router, authService, merchService, coinService, db, cfg.JWTSecret)

	// Запускаем сервер на указанном порту (из конфига)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
