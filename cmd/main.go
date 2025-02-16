package main

import (
	"avito_merchStore/internal/config"
	"avito_merchStore/internal/repository"
	"avito_merchStore/internal/service"
	"fmt"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	authService := service.NewAuthService(db, cfg.JWTSecret)

	// Тест логина
	token, err := authService.Login("irina", "mypassword")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешный вход, токен:", token)
}
