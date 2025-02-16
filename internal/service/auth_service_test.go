package service_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"avito_merchStore/internal/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Хеш, соответствующий "password123".
// Сгенерирован через bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost).
const bcryptHashPassword123 = "$2a$10$FW7n2Zk8/2BrkRjfEZmnuOKb0g1/QxSZGop0kCcLbdL1Tx5gXZd36"

func TestAuthService_Authenticate_NewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Нет пользователя alice
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, password_hash, coins FROM users WHERE username=$1")).
		WithArgs("alice").
		WillReturnError(sql.ErrNoRows)

	// Создаём нового пользователя (id=1) с 1000 монет
	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs("alice", sqlmock.AnyArg(), 1000).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	authSvc := service.NewAuthService(db, "testSecret")
	token, err := authSvc.Authenticate("alice", "password123")

	assert.NoError(t, err, "При регистрации нового пользователя ошибки быть не должно")
	assert.NotEmpty(t, token, "Токен не должен быть пустым для нового пользователя")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthService_Authenticate_ExistingUser_WrongPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Для charlie в БД лежит хеш от "password123",
	// но мы попробуем ввести "wrongpass"
	rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "coins"}).
		AddRow(3, "charlie", bcryptHashPassword123, 1000)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, password_hash, coins FROM users WHERE username=$1")).
		WithArgs("charlie").
		WillReturnRows(rows)

	authSvc := service.NewAuthService(db, "testSecret")
	token, err := authSvc.Authenticate("charlie", "wrongpass")

	assert.Error(t, err, "Пароль неверный, должна быть ошибка")
	assert.Contains(t, err.Error(), "неверные учетные данные")
	assert.Empty(t, token, "Токен не выдается при неверном пароле")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthService_Authenticate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Эмуляция системной ошибки БД
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, password_hash, coins FROM users WHERE username=$1")).
		WithArgs("dave").
		WillReturnError(errors.New("DB is down"))

	authSvc := service.NewAuthService(db, "testSecret")
	token, err := authSvc.Authenticate("dave", "whatever")

	assert.Error(t, err, "Должна быть ошибка от БД")
	assert.Empty(t, token, "Токен не выдаётся при ошибке БД")

	assert.NoError(t, mock.ExpectationsWereMet())
}
