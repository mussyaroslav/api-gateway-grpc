package models

import "github.com/google/uuid"

// AuthRequest предназначена для объединения данных, получаемых во время регистрации
type AuthRequest struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type AuthResponse struct {
	JWTToken string `json:"jwt_token"`
}

// TokenInfo содержит информацию, извлеченную из JWT токена
type TokenInfo struct {
	UserID  string   // ID пользователя (из поля sub)
	Email   string   // Email пользователя
	Roles   []string // Роли пользователя
	IsValid bool     // Флаг валидности токена
}

type User struct {
	UserId       uuid.UUID `db:"user_id" json:"user_id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	Roles        []string  `db:"roles" json:"roles"`
	PasswordHash string    `db:"password_hash" json:"-"` // Не включаем в JSON
}
