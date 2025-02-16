package e2e_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // драйвер для PostgreSQL
	"github.com/stretchr/testify/assert"

	"avito_merchStore/internal/config"
	"avito_merchStore/internal/repository"
	"avito_merchStore/internal/routes"
	"avito_merchStore/internal/service"
)

// AuthResponse описывает ответ от /api/auth
type AuthResponse struct {
	Token string `json:"token"`
}

// InfoResponse описывает ответ от /api/info
type InfoResponse struct {
	Coins     int `json:"coins"`
	Inventory []struct {
		Type     string `json:"type"`
		Quantity int    `json:"quantity"`
	} `json:"inventory"`
	CoinHistory map[string][]struct {
		FromUser string `json:"fromUser,omitempty"`
		ToUser   string `json:"toUser,omitempty"`
		Amount   int    `json:"amount"`
	} `json:"coinHistory"`
}

// setupTestServer подключается к PostgreSQL, используя переменные окружения,
// очищает таблицы и поднимает HTTP-сервер с зарегистрированными роутами.
func setupTestServer() (*httptest.Server, *sql.DB, error) {
	// Устанавливаем переменные окружения для тестового подключения:
	os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("DATABASE_PORT", "5432")
	os.Setenv("DATABASE_USER", "postgres")
	os.Setenv("DATABASE_PASSWORD", "password")
	os.Setenv("DATABASE_NAME", "shop")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("JWT_SECRET", "testSecret")

	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключаемся к PostgreSQL с помощью вашей функции
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Очистим таблицы, чтобы тесты работали с чистым состоянием.
	_, err = db.Exec("TRUNCATE TABLE transactions, purchases, users RESTART IDENTITY")
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка очистки таблиц: %w", err)
	}

	// Создаём сервисы
	authService := service.NewAuthService(db, cfg.JWTSecret)
	merchService := service.NewMerchService(db)
	coinService := service.NewCoinService(db)

	// Создаём роутер Gin и регистрируем маршруты.
	router := gin.Default()
	routes.RegisterRoutes(router, authService, merchService, coinService, db, cfg.JWTSecret)

	// Поднимаем тестовый HTTP-сервер
	ts := httptest.NewServer(router)
	return ts, db, nil
}

// loginHelper выполняет POST-запрос к /api/auth и возвращает JWT-токен.
func loginHelper(t *testing.T, baseURL, username, password string) string {
	body := fmt.Sprintf(`{"username": %q, "password": %q}`, username, password)
	resp, err := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err, "Ошибка при выполнении запроса /api/auth")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Ожидается успешный логин, статус 200")

	var ar AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&ar)
	assert.NoError(t, err, "Ошибка декодирования ответа /api/auth")
	assert.NotEmpty(t, ar.Token, "Токен не должен быть пустым после логина")
	return ar.Token
}

func TestE2E_BuyMerch(t *testing.T) {
	ts, db, err := setupTestServer()
	if err != nil {
		t.Fatalf("Не удалось запустить тестовый сервер: %v", err)
	}
	defer ts.Close()
	defer db.Close()

	// 1. Логинимся как alice через /api/auth
	token := loginHelper(t, ts.URL, "alice", "1234")

	// 2. Выполняем покупку мерча: покупаем "cup"
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/buy/cup", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Покупка мерча должна пройти успешно (статус 200)")

	// 3. Проверяем состояние пользователя через /api/info
	reqInfo, err := http.NewRequest(http.MethodGet, ts.URL+"/api/info", nil)
	assert.NoError(t, err)
	reqInfo.Header.Set("Authorization", "Bearer "+token)
	respInfo, err := http.DefaultClient.Do(reqInfo)
	assert.NoError(t, err)
	defer respInfo.Body.Close()

	var info InfoResponse
	err = json.NewDecoder(respInfo.Body).Decode(&info)
	assert.NoError(t, err, "Ошибка декодирования ответа /api/info")
	// При регистрации пользователь получает 1000 монет, а cup стоит 20
	assert.Equal(t, 980, info.Coins, "Баланс должен быть 1000 - 20 = 980")

	foundCup := false
	for _, item := range info.Inventory {
		if item.Type == "cup" && item.Quantity == 1 {
			foundCup = true
			break
		}
	}
	assert.True(t, foundCup, "Инвентарь должен содержать купленный мерч 'cup'")
}

func TestE2E_TransferCoins(t *testing.T) {
	ts, db, err := setupTestServer()
	if err != nil {
		t.Fatalf("Не удалось запустить тестовый сервер: %v", err)
	}
	defer ts.Close()
	defer db.Close()

	// 1. Логинимся под двумя пользователями через /api/auth
	tokenAlice := loginHelper(t, ts.URL, "alice", "1234")
	tokenBob := loginHelper(t, ts.URL, "bob", "5678")

	// 2. alice переводит bob 200 монет через /api/sendCoin
	transferBody := `{"toUser": "bob", "amount": 200}`
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/sendCoin", bytes.NewBufferString(transferBody))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenAlice)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Трансфер должен пройти успешно (статус 200)")

	// Немного подождем для завершения транзакций
	time.Sleep(100 * time.Millisecond)

	// 3. Проверяем баланс alice через /api/info (должно быть 800)
	reqInfoAlice, err := http.NewRequest(http.MethodGet, ts.URL+"/api/info", nil)
	assert.NoError(t, err)
	reqInfoAlice.Header.Set("Authorization", "Bearer "+tokenAlice)
	respInfoAlice, err := http.DefaultClient.Do(reqInfoAlice)
	assert.NoError(t, err)
	defer respInfoAlice.Body.Close()

	var infoAlice InfoResponse
	err = json.NewDecoder(respInfoAlice.Body).Decode(&infoAlice)
	assert.NoError(t, err)
	assert.Equal(t, 800, infoAlice.Coins, "У alice должно остаться 800 монет после перевода 200")

	// 4. Проверяем баланс bob через /api/info (должно быть 1200)
	reqInfoBob, err := http.NewRequest(http.MethodGet, ts.URL+"/api/info", nil)
	assert.NoError(t, err)
	reqInfoBob.Header.Set("Authorization", "Bearer "+tokenBob)
	respInfoBob, err := http.DefaultClient.Do(reqInfoBob)
	assert.NoError(t, err)
	defer respInfoBob.Body.Close()

	var infoBob InfoResponse
	err = json.NewDecoder(respInfoBob.Body).Decode(&infoBob)
	assert.NoError(t, err)
	assert.Equal(t, 1200, infoBob.Coins, "У bob должно быть 1200 монет после получения перевода")
}
