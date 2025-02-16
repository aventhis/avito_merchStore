// Package service реализует бизнес-логику приложения.
package service

import (
	"avito_merchStore/internal/models"
	"database/sql"
	"errors"
)

// CoinService предоставляет методы для операций с монетами.
type CoinService struct {
	db *sql.DB
}

// NewCoinService создаёт новый экземпляр CoinService.
func NewCoinService(db *sql.DB) *CoinService {
	return &CoinService{db: db}
}

// TransferCoins переводит монеты от одного пользователя к другому.
// fromUsername – имя отправителя (для записи в историю транзакций).
func (s *CoinService) TransferCoins(fromUserID int64, fromUsername, toUsername string, amount int) error {
	if amount <= 0 {
		return errors.New("сумма должна быть положительной")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем баланс отправителя
	var senderCoins int
	err = tx.QueryRow("SELECT coins FROM users WHERE id=$1 FOR UPDATE", fromUserID).Scan(&senderCoins)
	if err != nil {
		return err
	}
	if senderCoins < amount {
		return errors.New("недостаточно монет для перевода")
	}

	// Получаем ID получателя
	var receiverID int64
	err = tx.QueryRow("SELECT id FROM users WHERE username=$1 FOR UPDATE", toUsername).Scan(&receiverID)
	if err == sql.ErrNoRows {
		return errors.New("получатель не найден")
	} else if err != nil {
		return err
	}
	// Обновляем баланс отправителя
	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id=$2", amount, fromUserID)
	if err != nil {
		return err
	}

	// Обновляем баланс получателя
	_, err = tx.Exec("UPDATE users SET coins = coins + $1 WHERE id=$2", amount, receiverID)
	if err != nil {
		return err
	}

	// Регистрируем транзакцию отправки
	_, err = tx.Exec("INSERT INTO transactions (user_id, type, amount, counterpart) VALUES ($1, $2, $3, $4)",
		fromUserID, models.TransactionTypeTransfer, amount, toUsername)
	if err != nil {
		return err
	}

	// Регистрируем транзакцию получения
	_, err = tx.Exec("INSERT INTO transactions (user_id, type, amount, counterpart) VALUES ($1, $2, $3, $4)",
		receiverID, models.TransactionTypeTransfer, amount, fromUsername)
	if err != nil {
		return err
	}

	return tx.Commit()
}
