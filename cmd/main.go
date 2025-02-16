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

	// Тест логина (создаст пользователя, если его нет)
	_, err = authService.Login("alice", "password1")
	if err != nil {
		log.Fatal("Ошибка входа Alice:", err)
	}
	_, err = authService.Login("bob", "password2")
	if err != nil {
		log.Fatal("Ошибка входа Bob:", err)
	}

	// Получаем ID пользователей
	var aliceID, bobID int64
	err = db.QueryRow("SELECT id FROM users WHERE username=$1", "alice").Scan(&aliceID)
	if err != nil {
		log.Fatal("Ошибка получения ID Alice:", err)
	}
	err = db.QueryRow("SELECT id FROM users WHERE username=$1", "bob").Scan(&bobID)
	if err != nil {
		log.Fatal("Ошибка получения ID Bob:", err)
	}

	fmt.Println("ID Alice:", aliceID)
	fmt.Println("ID Bob:", bobID)

	// Создаём сервис для перевода монет
	coinService := service.NewCoinService(db)

	// Пробуем перевести 300 монет от Alice к Bob
	amount := 300
	err = coinService.TransferCoins(aliceID, "alice", "bob", amount)
	if err != nil {
		log.Fatalf("Ошибка перевода монет: %v", err)
	}

	fmt.Printf("Успешный перевод %d монет от Alice к Bob!\n", amount)

	// Проверяем балансы после перевода
	var aliceCoins, bobCoins int
	err = db.QueryRow("SELECT coins FROM users WHERE id=$1", aliceID).Scan(&aliceCoins)
	if err != nil {
		log.Fatal("Ошибка получения баланса Alice:", err)
	}
	err = db.QueryRow("SELECT coins FROM users WHERE id=$1", bobID).Scan(&bobCoins)
	if err != nil {
		log.Fatal("Ошибка получения баланса Bob:", err)
	}

	fmt.Printf("Баланс Alice: %d монет\n", aliceCoins)
	fmt.Printf("Баланс Bob: %d монет\n", bobCoins)
}
