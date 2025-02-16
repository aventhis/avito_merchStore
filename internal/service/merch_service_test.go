package service_test

import (
	"regexp"
	"testing"

	"avito_merchStore/internal/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMerchService_PurchaseMerch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Пользователь (id=1) имеет 30 монет. Покупает "cup" за 20.
	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT coins FROM users WHERE id=$1 FOR UPDATE")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(30))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET coins = coins - $1 WHERE id=$2")).
		WithArgs(20, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO purchases (user_id, item, quantity)
					  VALUES ($1, $2, 1)
					  ON CONFLICT (user_id, item) DO UPDATE SET quantity = purchases.quantity + 1`)).
		WithArgs(int64(1), "cup").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (user_id, type, amount, counterpart) VALUES ($1, $2, $3, $4)")).
		WithArgs(int64(1), "PURCHASE", 20, "cup").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	merchSvc := service.NewMerchService(db)
	err = merchSvc.PurchaseMerch(1, "cup")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMerchService_PurchaseMerch_NotEnoughCoins(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Пользователь (id=1) имеет 10 монет, "cup" стоит 20
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT coins FROM users WHERE id=$1 FOR UPDATE")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(10))

	mock.ExpectRollback()

	merchSvc := service.NewMerchService(db)
	err = merchSvc.PurchaseMerch(1, "cup")
	assert.EqualError(t, err, "недостаточно монет")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMerchService_PurchaseMerch_InvalidItem(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// В service.MerchList нет "unicorn"
	merchSvc := service.NewMerchService(db)
	err = merchSvc.PurchaseMerch(1, "unicorn")

	assert.EqualError(t, err, "неверное тип мерча")
}
