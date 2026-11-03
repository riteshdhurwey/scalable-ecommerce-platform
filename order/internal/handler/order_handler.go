package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/service"
)


type OrderHandler struct {

	service *service.OrderService
}

func InitOrderHandler(service *service.OrderService) *OrderHandler{
	return &OrderHandler{service: service}
}

func(h *OrderHandler)CreateOrder(ctx *gin.Context){
    token,_ := ctx.Cookie("authorization")

    if token == "" {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": "authorization header missing",
        })
        return
    }

    value, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": "userID missing",
        })
        return
    }

    userID, ok := value.(int64)
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": "invalid userID",
        })
        return
    }

	outOfStockProducts,availableProducts,totalAmount,err := h.service.CheckOut(ctx.Request.Context(),userID,token)

	if err != nil{
		if errors.Is(err,service.ErrOutOfStock){
			ctx.JSON(http.StatusConflict,gin.H{
				"error":err.Error(),"out of stock":outOfStockProducts,
			})
			return
		}
		if errors.Is(err,service.ErrEmptyCart){
			ctx.JSON(http.StatusConflict,gin.H{
				"error":err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK,gin.H{
		"Order Placed":availableProducts,
		"total Amount":totalAmount,
	})
}