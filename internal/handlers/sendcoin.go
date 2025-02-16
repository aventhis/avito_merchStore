package handlers

import (
	"net/http"

	"avito_merchStore/internal/service"
	"github.com/gin-gonic/gin"
)

type SendCoinHandler struct {
	coinService *service.CoinService
}

func NewSendCoinHandler(coinService *service.CoinService) *SendCoinHandler {
	return &SendCoinHandler{coinService: coinService}
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

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
