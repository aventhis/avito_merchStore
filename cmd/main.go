package main

import (
	"avito_merchStore/internal/config"
	"avito_merchStore/internal/repository"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к бд", err)
	}
	defer db.Close()

}
