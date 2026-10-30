package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/dto"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func InitUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(ctx *gin.Context) {

	var req dto.RegisterRequest

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, userErr := h.service.Create(ctx, &req)

	if userErr != nil {
		switch {
		case errors.Is(userErr, service.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": userErr.Error(),
			})

		case errors.Is(userErr, service.ErrEmailRegistered):
			ctx.JSON(http.StatusConflict, gin.H{
				"error": userErr.Error(),
			})

		case errors.Is(userErr, service.ErrUserIDRegistered):
			ctx.JSON(http.StatusConflict, gin.H{
				"error": userErr.Error(),
			})

		case errors.Is(userErr, service.ErrUserNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": userErr.Error(),
			})

		case errors.Is(userErr, service.ErrInvalidCredentials):
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": userErr.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var request dto.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	authToken, err := h.service.Login(ctx, request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})

		case errors.Is(err, service.ErrInvalidCredentials):
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	if authToken == nil || *authToken == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate access token",
		})
		return
	}

	ctx.SetCookie(
		"authorization",
		*authToken,
		900,
		"/",
		"",
		true, // Secure: HTTPS only
		true, // HttpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
	})
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userName := ctx.GetString("username")

	if userName == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "username missing",
		})
		return
	}

	user, err := h.service.GetProfile(ctx, userName)

	if err != nil {
		if errors.Is(err,service.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}


func (h *UserHandler) DeleteUser(ctx *gin.Context){

	userName := ctx.GetString("username")

	_,err := h.service.DeleteUser(ctx,userName)

	if err != nil {
		if errors.Is(err,service.ErrUserNotFound){
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

    ctx.JSON(http.StatusAccepted,gin.H{"response":"user deleted"})

}