package handler

import (
	"context"
	"errors"
	"log"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/service"
	pb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGrpcHandler struct {
	pb.UnimplementedOrderServiceServer

	orderService *service.OrderService
}

func InitOrderGrpcHandler(orderservice *service.OrderService) *OrderGrpcHandler{
	return &OrderGrpcHandler{orderService: orderservice}
}

func (h *OrderGrpcHandler) UpdateExpiredOrderStatus(ctx context.Context,req *pb.GetOrderStatusRequest)(*pb.GetOrderStatusResponse,error){

	log.Println("worker request recieved")

	response,err := h.orderService.UpdateExpiredOrderStatus(ctx,req)

	if err != nil {
		if errors.Is(err,service.ErrInvalidOrderId) {
			return nil,status.Error(codes.InvalidArgument,err.Error())
		}

		return nil,status.Error(codes.Internal,err.Error())

	}

	log.Println("worker request done")

	return response,nil
}