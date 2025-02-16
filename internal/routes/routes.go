// Package routes определяет маршруты HTTP-сервера.
package routes

import (
	"avito_merchStore/internal/handlers"
	"avito_merchStore/internal/middleware"
	"avito_merchStore/internal/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes регистрирует публичные и защищённые маршруты для приложения.
func RegisterRoutes(router *gin.Engine, authService *service.AuthService, merchService *service.MerchService, coinService *service.CoinService, db *sql.DB, jwtSecret string) {
	// Публичный эндпоинт для аутентификации
	authHandler := handlers.NewAuthHandler(authService)
	router.POST("/api/auth", authHandler.Login)

	// Группа защищенных эндпоинтов с JWT-мидлваром
	protected := router.Group("/", middleware.JWTAuthMiddleware(jwtSecret))
	{
		infoHandler := handlers.NewInfoHandler(db)
		protected.GET("/api/info", infoHandler.GetInfo)

		buyHandler := handlers.NewBuyHandler(merchService)
		protected.GET("/api/buy/:item", buyHandler.BuyMerch)

		sendCoinHandler := handlers.NewSendCoinHandler(coinService)
		protected.POST("/api/sendCoin", sendCoinHandler.SendCoin)
	}
}
