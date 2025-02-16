package service

import (
	"avito_merchStore/internal/models"
	"database/sql"
	"errors"
)

type MerchService struct {
	db *sql.DB
}

func NewMerchService(db *sql.DB) *MerchService {
	return &MerchService{db: db}
}

func (s *MerchService) PurchaseMerch(userID int64, item string) error {
	var price int
	found := false
	for _, m := range models.MerchList {
		if m.Name == item {
			price = m.Price
			found = true
			break
		}
	}
	if !found {
		return errors.New("неверное тип мерча")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем баланс пользователя
	var coins int
	err = tx.QueryRow("SELECT coins FROM users WHERE id=$1 FOR UPDATE", userID).Scan(&coins)
	if err != nil {
		return errors.New("ошибка получения баланса пользователя")
	}
	if coins < price {
		return errors.New("недостаточно монет")
	}

	// Списываем монеты
	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id=$2", price, userID)
	if err != nil {
		return err
	}
	// Обновляем инвентарь (таблица purchases)
	_, err = tx.Exec(`INSERT INTO purchases (user_id, item, quantity)
					  VALUES ($1, $2, 1)
					  ON CONFLICT (user_id, item) DO UPDATE SET quantity = purchases.quantity + 1`,
		userID, item)
	if err != nil {
		return err
	}

	// Регистрируем транзакцию покупки
	_, err = tx.Exec("INSERT INTO transactions (user_id, type, amount, counterpart) VALUES ($1, $2, $3, $4)",
		userID, models.TransactionTypePurchase, price, item)
	if err != nil {
		return err
	}

	return tx.Commit()
}
