package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	db *sql.DB
}

func NewInfoHandler(db *sql.DB) *InfoHandler {
	return &InfoHandler{db: db}
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinTransaction struct {
	// Для полученных транзакций – от кого получены монеты
	FromUser string `json:"fromUser,omitempty"`
	// Для отправленных транзакций – кому отправлены монеты
	ToUser string `json:"toUser,omitempty"`
	Amount int    `json:"amount"`
}

type InfoResponse struct {
	Coins       int                          `json:"coins"`
	Inventory   []InventoryItem              `json:"inventory"`
	CoinHistory map[string][]CoinTransaction `json:"coinHistory"`
}

func (h *InfoHandler) GetInfo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// Получаем баланс пользователя
	var coins int
	err := h.db.QueryRow("SELECT coins FROM users WHERE id=$1", userID).Scan(&coins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	// Получаем инвентарь (например, покупки мерча)
	rowsInv, err := h.db.Query("SELECT item, quantity FROM purchases WHERE user_id=$1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}
	defer rowsInv.Close()
	inventory := []InventoryItem{}
	for rowsInv.Next() {
		var item InventoryItem
		if err := rowsInv.Scan(&item.Type, &item.Quantity); err != nil {
			continue
		}
		inventory = append(inventory, item)
	}

	// Получаем историю транзакций (для переводов)
	rowsTx, err := h.db.Query("SELECT type, amount, counterpart FROM transactions WHERE user_id=$1 AND type='transfer'", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}
	defer rowsTx.Close()

	var sent []CoinTransaction
	var received []CoinTransaction

	for rowsTx.Next() {
		var txType string
		var amount int
		var counterpart string
		if err := rowsTx.Scan(&txType, &amount, &counterpart); err != nil {
			continue
		}
		// Если сумма отрицательная – это исходящая транзакция (sent)
		if amount < 0 {
			sent = append(sent, CoinTransaction{
				ToUser: counterpart,
				Amount: -amount, // берем абсолютное значение
			})
		} else {
			// Сумма положительная – входящая транзакция (received)
			received = append(received, CoinTransaction{
				FromUser: counterpart,
				Amount:   amount,
			})
		}
	}

	coinHistory := map[string][]CoinTransaction{
		"received": received,
		"sent":     sent,
	}

	c.JSON(http.StatusOK, InfoResponse{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	})
}
