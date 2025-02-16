package service_test

import (
	"database/sql"
	"regexp"
	"testing"

	"avito_merchStore/internal/models"
	"avito_merchStore/internal/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCoinService_TransferCoins_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	// Отправитель (id=1) имеет 50 монет
	mock.ExpectQuery(regexp.QuoteMeta("SELECT coins FROM users WHERE id=$1 FOR UPDATE")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(50))

	// Получатель "bob" имеет id=2
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE username=$1 FOR UPDATE")).
		WithArgs("bob").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	// Списываем монеты у отправителя
	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET coins = coins - $1 WHERE id=$2")).
		WithArgs(10, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Добавляем монеты получателю
	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET coins = coins + $1 WHERE id=$2")).
		WithArgs(10, int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Запись транзакции у отправителя
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (user_id, type, amount, counterpart) VALUES ($1, $2, $3, $4)")).
		WithArgs(int64(1), models.TransactionTypeTransfer, 10, "bob").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Запись транзакции у получателя
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (user_id, type, amount, counterpart) VALUES ($1, $2, $3, $4)")).
		WithArgs(int64(2), models.TransactionTypeTransfer, 10, "alice").
		WillReturnResult(sqlmock.NewResult(2, 1))

	mock.ExpectCommit()

	coinSvc := service.NewCoinService(db)
	err = coinSvc.TransferCoins(1, "alice", "bob", 10)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoinService_TransferCoins_NotEnoughCoins(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	// Отправитель (id=1) имеет 5 монет
	mock.ExpectQuery(regexp.QuoteMeta("SELECT coins FROM users WHERE id=$1 FOR UPDATE")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(5))

	// Транзакция должна быть откатана
	mock.ExpectRollback()

	coinSvc := service.NewCoinService(db)
	err = coinSvc.TransferCoins(1, "alice", "bob", 10)
	assert.EqualError(t, err, "недостаточно монет для перевода")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoinService_TransferCoins_RecipientNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	// sender has 50 coins
	mock.ExpectQuery(regexp.QuoteMeta("SELECT coins FROM users WHERE id=$1 FOR UPDATE")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(50))

	// получатель не найден
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE username=$1 FOR UPDATE")).
		WithArgs("bob").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectRollback()

	coinSvc := service.NewCoinService(db)
	err = coinSvc.TransferCoins(1, "alice", "bob", 10)
	assert.EqualError(t, err, "получатель не найден")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoinService_TransferCoins_AmountMustBePositive(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	coinSvc := service.NewCoinService(db)
	err = coinSvc.TransferCoins(1, "alice", "bob", 0)
	assert.EqualError(t, err, "сумма должна быть положительной")
}
