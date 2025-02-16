// Package models содержит определения структур данных, используемых в приложении.
package models

// User представляет пользователя приложения.
type User struct {
	ID           int64  `db:"id"`            // Уникальный идентификатор пользователя
	Username     string `db:"username"`      // Имя пользователя
	PasswordHash string `db:"password_hash"` // Хэш пароля пользователя
	Coins        int    `db:"coins"`         // Баланс монет пользователя
}
