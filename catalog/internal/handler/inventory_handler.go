package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/service"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func InitInventoryHandler(service *service.InventoryService) *InventoryHandler{
	return &InventoryHandler{service: service}
}


func (h *InventoryHandler) ReserveInventory(ctx *gin.Context){

	var reserveReq models.ReservationRequest

	err := ctx.ShouldBindJSON(&reserveReq)

	if err != nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	err = h.service.ReserveInventory(ctx,reserveReq)

	if err != nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})

		return
	}

	ctx.JSON(http.StatusOK,gin.H{"response":"ok"})
}