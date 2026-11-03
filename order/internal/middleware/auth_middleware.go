package middleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/jwt"
)



func AuthMiddleware(ctx *gin.Context){
	tokenString,err := ctx.Cookie("authorization")

	if err != nil || tokenString == ""{
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
			"error":"authentication required",
		})
		return
	}

	claim,err := jwt.ValidateToken(tokenString)

	if err != nil {
		if errors.Is(err,jwt.ErrExpiredToken){
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"error":"token expired",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
			"error":"invalid token",
		})
		return
	}

	subject,ok := claim["sub"].(string)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
			"error":"invalid token claims",
		})
		return
	}

	userID,err := strconv.ParseInt(subject,10,64)

	if err != nil || userID <= 0{
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
			"error":"invalid token claims",
		})

		return
	}

	ctx.Set("userID",userID)
	ctx.Next()

}