package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/jwt"
)

func AuthMiddleware(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("authorization")
	if err != nil || tokenString == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		if errors.Is(err, jwt.ErrExpiredToken) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	userName, ok := claims["sub"].(string)
	if !ok || userName == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token claims",
		})
		return
	}

	ctx.Set("username", userName)
	ctx.Next()
}
