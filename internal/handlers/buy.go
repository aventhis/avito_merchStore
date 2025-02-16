package handlers

import (
	"net/http"

	"avito_merchStore/internal/service"
	"github.com/gin-gonic/gin"
)

type BuyHandler struct {
	merchService *service.MerchService
}

func NewBuyHandler(merchService *service.MerchService) *BuyHandler {
	return &BuyHandler{merchService: merchService}
}

func (h *BuyHandler) BuyMerch(c *gin.Context) {
	userID := c.GetInt64("user_id")
	item := c.Param("item")
	if item == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Параметр item обязателен"})
		return
	}
	err := h.merchService.PurchaseMerch(userID, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Покупка успешна"})
}
