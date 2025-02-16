package models

import "time"

type TransactionType string

const (
	TransactionTypePurchase TransactionType = "PURCHASE"
	TransactionTypeCredit   TransactionType = "TRANSFER"
)

type Transaction struct {
	ID          int64           `db:"id"`
	UserID      int64           `db:"user_id"`
	Type        TransactionType `db:"type"`
	Amount      int             `db:"amount"`
	Counterpart string          `db:"counterpart"` // для transfer – имя контрагента, для покупки – название товара
	CreatedAt   time.Time       `db:"created_at"`
}
