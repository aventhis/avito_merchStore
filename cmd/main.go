package main

import (
	"avito_merchStore/internal/config"
	"avito_merchStore/internal/repository"
	"avito_merchStore/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
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

}
