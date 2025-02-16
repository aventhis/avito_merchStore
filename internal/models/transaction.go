// Package models содержит определения структур данных для транзакций и мерча.
package models

import "time"

// TransactionType определяет тип транзакции.
type TransactionType string

const (
	// TransactionTypePurchase обозначает транзакцию покупки мерча.
	TransactionTypePurchase TransactionType = "PURCHASE"
	// TransactionTypeTransfer обозначает транзакцию перевода монет.
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

// Transaction представляет запись о транзакции в системе.
type Transaction struct {
	ID          int64           `db:"id"`          // Уникальный идентификатор транзакции
	UserID      int64           `db:"user_id"`     // Идентификатор пользователя, совершившего транзакцию
	Type        TransactionType `db:"type"`        // Тип транзакции (PURCHASE или TRANSFER)
	Amount      int             `db:"amount"`      // Сумма транзакции
	Counterpart string          `db:"counterpart"` // Для transfer – имя контрагента, для покупки – название товара
	CreatedAt   time.Time       `db:"created_at"`  // Время создания транзакции
}
