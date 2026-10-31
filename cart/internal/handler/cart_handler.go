package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/service"
)


type CartHandler struct {
	service *service.CartSevice
}


func InitCartHandler(service *service.CartSevice) *CartHandler {
	return &CartHandler{service: service}
}


func (h *CartHandler)AddToCart(ctx *gin.Context){

    value,exists := ctx.Get("userID")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "userID missing",
		})
		return
	}
		userID, ok := value.(int64)

	if !ok || userID <= 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid userID",
			})
			return
		}

	var product models.Cart

	err := ctx.ShouldBindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"invalid parameters"})
		return
	}

	err = h.service.AddToCart(ctx,int64(userID),product)

	if  err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	ctx.JSON(http.StatusOK,gin.H{"responce":"added to cart"})
}


func (h *CartHandler) GetCart(ctx *gin.Context){

	value,exists := ctx.Get("userID")

	if !exists {

		ctx.JSON(http.StatusUnauthorized,gin.H{"error":"userID missing"})
		return
	}

	userID,ok := value.(int64)

	if !ok || userID <= 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid userID",
			})
			return
	}

	products,err := h.service.GetCart(ctx,userID)


	if err != nil {

		if errors.Is(err,service.ErrCartIDNotFound) {
			ctx.JSON(http.StatusOK,gin.H{"response":"cart id not found"})
			return
		}
		
		if errors.Is(err,service.ErrNoItemFound) {
			ctx.JSON(http.StatusOK,gin.H{"response":"cart is empty"})
			return
		}

		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

	ctx.JSON(http.StatusOK,products)
}