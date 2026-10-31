package dto

import "time"

type RegisterRequest struct {
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UserID   *int64  `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	UserName  string    `json:"user_name"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	PasswordHash string `json:"password_hash"`
}
