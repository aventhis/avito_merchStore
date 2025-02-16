// Package config предоставляет функции для загрузки конфигурационных параметров
// из переменных окружения с заданными значениями по умолчанию.
package config

import "os"

// Config содержит настройки подключения к базе данных, настройки сервера и секрет для JWT.
type Config struct {
	DBHost     string // Хост базы данных (например, localhost)
	DBPort     string // Порт базы данных (например, 5432)
	DBUser     string // Имя пользователя базы данных (например, postgres)
	DBPassword string // Пароль пользователя базы данных
	DBName     string // Название базы данных (например, shop)
	ServerPort string // Порт, на котором запущен сервер (например, 8080)
	JWTSecret  string // Секрет для подписывания JWT-токенов
}

// LoadConfig загружает конфигурацию из переменных окружения. Если переменная не задана,
// используется значение по умолчанию.
func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DATABASE_HOST", "localhost"),
		DBPort:     getEnv("DATABASE_PORT", "5432"),
		DBUser:     getEnv("DATABASE_USER", "postgres"),
		DBPassword: getEnv("DATABASE_PASSWORD", "password"),
		DBName:     getEnv("DATABASE_NAME", "shop"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
