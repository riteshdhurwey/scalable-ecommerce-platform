package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

	orderProduct,err := h.service.AddToCart(ctx,int64(userID),product)

	if  err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	ctx.JSON(http.StatusOK,gin.H{"responce":"added to cart","product":orderProduct})
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
		
		if errors.Is(err,service.ErrEmptyCart) {
			ctx.JSON(http.StatusOK,gin.H{"response":"cart is empty"})
			return
		}

		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

	ctx.JSON(http.StatusOK,products)
}

func (h *CartHandler) UpdateCart(ctx *gin.Context){
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

	var req models.Cart

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"invalid update parameters"})
		return
	}

_,err = h.service.UpdateCart(ctx,userID,req)

	if err != nil {

		switch {
		case errors.Is(err,service.ErrCartIDNotFound) :
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
		
		case errors.Is(err,service.ErrProductNotFound) :
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
		
		default :
			
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		
		}
		return
	}


	ctx.JSON(http.StatusCreated,gin.H{"response":"cart updated"})
}


func (h *CartHandler) DeleteItem(ctx *gin.Context){
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

	idStr := ctx.Param("id")

	productID,err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"invalid product id"})
		return
	}

	err = h.service.DeleteItem(ctx,userID,int64(productID))

	if err != nil {

		switch {
		case errors.Is(err,service.ErrCartIDNotFound) :
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
		
		case errors.Is(err,service.ErrProductNotFound) :
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
		
		default :
			
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		
		}
		return
	}

	ctx.JSON(http.StatusOK,gin.H{"response":"item deleted"})


}


func (h *CartHandler) CheckOut(ctx *gin.Context){

	value,exists := ctx.Get("userID")
	fmt.Printf("userID type: %T\n", value)

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

	totalAmount,products,outOfStockProducts,err := h.service.CheckOut(ctx,userID)


	if err != nil {

		if errors.Is(err,service.ErrCartIDNotFound) {
			ctx.JSON(http.StatusConflict,gin.H{"response":"cart id not found"})
			return
		}
		
		if errors.Is(err,service.ErrEmptyCart) {
			ctx.JSON(http.StatusOK,gin.H{"response":"cart is empty"})
			return
		}

		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}


	ctx.JSON(http.StatusOK,gin.H{"totalAmount":totalAmount,"available":products,"out of stock":outOfStockProducts})
}

func (h *CartHandler) ClearCart(c *gin.Context) {
    value, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    userID := value.(int64)

    err := h.service.ClearCart(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "cart cleared successfully",
    })
}
