// Package handlers содержит HTTP-обработчики для операций с монетами.
package handlers

import (
	"net/http"

	"avito_merchStore/internal/service"
	"github.com/gin-gonic/gin"
)

// SendCoinHandler обрабатывает запросы на перевод монет.
type SendCoinHandler struct {
	coinService *service.CoinService
}

// NewSendCoinHandler создаёт новый экземпляр SendCoinHandler.
func NewSendCoinHandler(coinService *service.CoinService) *SendCoinHandler {
	return &SendCoinHandler{coinService: coinService}
}

// SendCoinRequest представляет тело запроса на перевод монет.
type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

// SendCoin обрабатывает запрос на перевод монет.
func (h *SendCoinHandler) SendCoin(c *gin.Context) {
	userID := c.GetInt64("user_id")
	username := c.GetString("username")
	var req SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Неверное тело запроса"})
		return
	}
	if req.ToUser == "" || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Неверные параметры запроса"})
		return
	}
	err := h.coinService.TransferCoins(userID, username, req.ToUser, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Перевод успешен"})
}
