// Package handlers предоставляет HTTP-обработчики для аутентификации пользователей.
package handlers

import (
	"avito_merchStore/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthHandler обрабатывает запросы, связанные с аутентификацией.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler создаёт новый AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// AuthRequest представляет тело запроса для аутентификации пользователя.
type AuthRequest struct {
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Пароль пользователя
}

// AuthResponse представляет ответ на успешную аутентификацию.
type AuthResponse struct {
	Token string `json:"token"` // JWT-токен
}

// Login обрабатывает HTTP-запрос для авторизации пользователя.
// При успешной аутентификации возвращается JWT-токен.
func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}
	token, err := h.authService.Authenticate(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}
