package handler

import (
	"context"
	"log"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/service"
	pb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/order"
)

type OrderGrpcHandler struct {
	pb.UnimplementedOrderServiceServer

	orderService *service.OrderService
}

func InitOrderGrpcHandler(orderservice *service.OrderService) *OrderGrpcHandler{
	return &OrderGrpcHandler{orderService: orderservice}
}

func (h *OrderGrpcHandler) GetOrderStatus(ctx context.Context,req *pb.GetOrderStatusRequest)(*pb.GetOrderStatusResponse,error){

	log.Println("worker request recieved")

	OrderExpired := make([]*pb.OrderState,0,len(req.OrderIDs))

	for _,order := range req.OrderIDs {

		expired,err := h.orderService.IsOrderExpired(ctx,order.OrderID)

		if err != nil {
			return nil,err
		}


		response := pb.OrderState{OrderID: order.OrderID,OrderExpired : expired}

		OrderExpired = append(OrderExpired, &response)
		
	}

	response := &pb.GetOrderStatusResponse{OrderStates: OrderExpired}
	log.Println("worker request done")
	return response,nil
}