package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"time"

	cartClient "github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/client/cart"
	inventoryClient "github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/client/inventory"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/repository"
	"google.golang.org/grpc/metadata"
	pb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/order"

)

var (
	ErrOutOfStock = errors.New("some products are out of stock")
	ErrEmptyCart  = errors.New("Cart is empty")
	ErrAuthCredentials = errors.New("authorization credentials are missing")
	ErrInvalidOrderId  = errors.New("invalid order id")
)


type OrderService struct {
	repo *repository.OrderRepo
	inventoryClient *inventoryClient.Inventory
	cartClient *cartClient.Cart
}

func InitOrderService(repo *repository.OrderRepo,client *inventoryClient.Inventory,cartclient *cartClient.Cart) *OrderService {
	return &OrderService{repo: repo,inventoryClient: client,cartClient: cartclient}
}

func (s *OrderService) CheckOut(ctx context.Context,userID int64,token string)(*[]models.OutOfStockProducts,*[]models.OrderItems,float32,error){

	if userID <= 0 || token == ""{
		return nil,nil,0,ErrAuthCredentials
	}

	md := metadata.Pairs(
		  "authorization","Bearer" + token,	
	)

	ctx = metadata.NewOutgoingContext(ctx,md)

	response,err := s.cartClient.GetCartItems(ctx,userID)

	if err != nil {
		return nil,nil,0,err
	}

	if response.Products == nil {
		return nil,nil,0,ErrEmptyCart
	}

	order := models.Order{UserID: userID,TotalAmount: response.TotalAmount,PaymentExpAt: time.Now().Add(10 * time.Minute)}

	orderID,err := s.repo.CreateOrder(ctx,order)

	if err != nil {
		return nil,nil,0,err
	}

	if orderID == nil {
		return nil,nil,0,fmt.Errorf("failed to create order")
	}

	var orderReq models.ReservationRequest
	orderReq.OrderID = *orderID

	var items []models.ReservationItem
	
	AvailableProducts := response.Products

	for _,p := range AvailableProducts {
		var item models.ReservationItem

		item.ProductID = p.ID
		item.Quantity = int(p.Quantity)

		items = append(items, item)
	}

	orderReq.Items = items

	orderReserved,err := s.inventoryClient.ReverseInventory(ctx,orderReq)

	if err != nil {
    return nil, nil, 0, fmt.Errorf("failed to reserve inventory: %w", err)
    }

	if !orderReserved.Done{
		return nil,nil,0,fmt.Errorf("failed to reserve order")
	} 

	var products[]models.OrderItems

	for _,p := range AvailableProducts {

		var product models.OrderItems

		product.ProductID = p.ID
		product.ProductName = p.Name
		product.Quantity = int(p.Quantity)
		product.UnitPrice = p.Price
		product.SubTotal = p.Price * float64(p.Quantity)
		products = append(products, product)

	}

	err = s.repo.CreateOrderItems(ctx,*orderID,products)

	if err != nil {
		return nil,nil,0,err
	}

	log.Print("order placed")

	cartCleared,err := s.cartClient.ClearCart(ctx,userID)

	if err != nil || !cartCleared.Done {
		return nil,nil,0,fmt.Errorf("cart not cleared %w",err)
	}


	return nil,&products,float32(order.TotalAmount),nil

	
}


func (s *OrderService)GetOrderStatus(ctx context.Context,orderID int64)(models.OrderStatusResponse,error){

	if orderID <= 0 {
		return models.OrderStatusResponse{},ErrInvalidOrderId
	}

	return s.repo.GetOrderStatus(ctx,orderID)
}


func (s *OrderService) IsOrderExpired(ctx context.Context,orderID int64)(bool,error){

	if orderID <= 0 {
		return false,ErrInvalidOrderId
	}

	orderStatus,err := s.repo.GetOrderStatus(ctx,orderID)

	if err != nil {
		return false,err
	}

	if (orderStatus.PaymentExpAt.IsZero() && !orderStatus.PaymentExpAt.Add(10 * time.Minute).After(time.Now())) && orderStatus.Status == "PENDING" {
		return true,nil
	}

	return false,nil
}


func (s *OrderService) UpdateExpiredOrderStatus(ctx context.Context, req *pb.GetOrderStatusRequest)(*pb.GetOrderStatusResponse,error){

	if req == nil {
		return nil,ErrInvalidOrderId
	}

	var reqOrderIDs []int64

	for _,order := range req.OrderIDs {

		reqOrderIDs = append(reqOrderIDs, order.OrderID)
	}

	updated,updatedOrderIDs,err := s.repo.UpdateExpiredOrderStatus(ctx,reqOrderIDs)

	if err != nil {
		return nil,err
	}

	if !updated {
		
		return &pb.GetOrderStatusResponse{OrderIDs: nil,StatusUpdated: false},nil
	}

	orderIDs := make([]*pb.OrderID,0,len(updatedOrderIDs))

	for _,order := range updatedOrderIDs {

		orderID := &pb.OrderID{OrderID: order}

		orderIDs = append(orderIDs, orderID)
	}

	response := &pb.GetOrderStatusResponse{OrderIDs: orderIDs,StatusUpdated: true}

	return response,nil

}