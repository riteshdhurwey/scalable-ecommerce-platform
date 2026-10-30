package service

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid user input")
	ErrEmailRegistered    = errors.New("email already registered")
	ErrUserIDRegistered   = errors.New("User name already Exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInternalDatabase   = errors.New("Internal Database Error")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token expired")
)
