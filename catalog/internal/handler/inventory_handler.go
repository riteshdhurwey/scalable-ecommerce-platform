package handler

import (
	"context"
	"log"

	pb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/catalog"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/service"
)

type InventoryHandler struct {
	service *service.InventoryService
	pb.UnimplementedCatalogServiceServer
}

func InitInventoryHandler(service *service.InventoryService) *InventoryHandler{
	return &InventoryHandler{service: service}
}


func (h *CatalogHandler) ReserveInventory(ctx context.Context,req *pb.ReserveInventoryRequest)(*pb.ReserveInventoryResponse,error){

	var reserveReq models.ReservationRequest

	reserveReq.OrderID = req.OrderId

	var items []models.ReservationItem 

	for _,item := range req.ReservationItems {

		var request models.ReservationItem

		request.ProductID = item.ProductId
		request.Quantity = int(item.Quantity)

		items = append(items, request)
	}

    log.Printf("Received inventory reservation request: order_id=%d", req.OrderId)
	reserveReq.Items = items

	err := h.inventoryService.ReserveInventory(ctx,reserveReq)

	if err != nil {
		return &pb.ReserveInventoryResponse{Done: false},err
	}

	return &pb.ReserveInventoryResponse{Done: true},nil
}