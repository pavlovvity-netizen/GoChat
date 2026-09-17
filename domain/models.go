package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("сообщение не найдено")

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Status string `json:"status"`
	Token  string `json:"token,omitempty"` // JWT-токен или идентификатор сессии
	User   User   `json:"user"`
}

type Message struct {
	ID       int       `json:"id"`
	Text     string    `json:"text"`
	PostTime time.Time `json:"post_time"`
	UserID   int       `json:"user_id"`
	Username string    `json:"user_name,omitempty"`
}

type Response struct {
	Status  string  `json:"status"`
	Message Message `json:"message"`
}

type Request struct {
	Text   string `json:"text"`
	UserID int    `json:"user_id"`
}
