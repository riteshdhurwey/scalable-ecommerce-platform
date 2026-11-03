package jwt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/dto"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

func GenerateAccessToken(req dto.LoginRequest) (string, error) {
	_ = godotenv.Load()

	secret := os.Getenv("JWT_SECRET")
	if secret== "" {
		return "", errors.New("JWT secret is not configured")
	}

	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(*req.UserID, 10),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	_ : godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT secret is not configured")
	}

	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
