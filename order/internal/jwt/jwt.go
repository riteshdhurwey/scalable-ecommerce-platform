package jwt

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
_:
	godotenv.Load()

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return nil, errors.New("JWT secret is not configured")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil,ErrExpiredToken
		}

		return nil,ErrInvalidToken
	}

	claim,ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil,ErrInvalidToken
	}

	return claim,nil
}
