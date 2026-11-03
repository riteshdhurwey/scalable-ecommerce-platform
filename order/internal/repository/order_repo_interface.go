package repository

import (
	"context"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/models"
)

type OrderRepository interface {

	CreateOrder(ctx context.Context,userID int64,req models.Order)(*int64,error)
	CreateOrderItems(ctx context.Context,OrderID int64,req models.OrderItems)(error)
}