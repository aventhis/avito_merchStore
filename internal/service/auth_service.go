package service

import (
	"avito_merchStore/internal/models"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	db        *sql.DB
	JWTSecret string
}

func NewAuthService(db *sql.DB, JWTSecret string) *AuthService {
	return &AuthService{
		db:        db,
		JWTSecret: JWTSecret,
	}
}

func (s *AuthService) Login(username string, password string) (string, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, username, password_hash, coins FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins)
	if err == sql.ErrNoRows {
		// Если пользователь не найден, создаем его автоматически
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		// Создаем пользователя с начальным балансом 1000 монет
		err = s.db.QueryRow("INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3) RETURNING id", username, string(hash), 1000).Scan(&user.ID)
		if err != nil {
			return "", err
		}
		user.Username = username
		user.PasswordHash = string(hash)
		user.Coins = 1000
	} else if err != nil {
		return "", err
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			return "", errors.New("неверные учетные данные")
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
