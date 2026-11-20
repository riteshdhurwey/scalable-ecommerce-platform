package handler

import (
	"context"
	"errors"


	pb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/cart"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartGrpcHandler struct {

	pb.UnimplementedCartServiceServer
	cartService *service.CartSevice
}


func InitCartGrpcHandler(cartservice *service.CartSevice) *CartGrpcHandler{
	return	&CartGrpcHandler{cartService: cartservice}
}

func (h *CartGrpcHandler) GetCart(ctx context.Context, req *pb.GetCartRequest)(*pb.GetCartResponse,error){


   totalAmount, availableProducts, outOfStockProducts, err :=
        h.cartService.CheckOut(ctx, req.UserID)

    if err != nil {
        if errors.Is(err, service.ErrCartIDNotFound) {
            return nil, status.Error(codes.NotFound, err.Error())
        }

        if errors.Is(err, service.ErrEmptyCart) {
            // An empty cart is a valid business state.
            return &pb.GetCartResponse{
                TotalAmount:        0,
                Products:           []*pb.Product{},
                OutOfStockProducts: []*pb.OutOfStockProduct{},
            }, nil
        }

        return nil, status.Error(codes.Internal, err.Error())
    }

    // Convert available products to Protobuf products.
    products := make([]*pb.Product, 0, len(*availableProducts))

    for _, p := range *availableProducts {
        product := &pb.Product{
            ID:       p.ProductID,
            Name:     p.Name,
            Price:    p.Price,
            Quantity: int32(p.StockQuantity),
        }

        products = append(products, product)
    }

    // Convert out-of-stock products to Protobuf products.
    outOfStock := make([]*pb.OutOfStockProduct, 0, len(outOfStockProducts))

    for _, p := range outOfStockProducts {
        product := &pb.OutOfStockProduct{
            ID:       p.ProductID,
            Name:     p.Name,
            AvailableQuantity: int32(p.AvailableStock),
            OrderedQuantity: int32(p.OrderedStock),
        }

        outOfStock = append(outOfStock, product)
    }

    response := &pb.GetCartResponse{
        TotalAmount:        float64(totalAmount),
        Products:           products,
        OutOfStockProducts: outOfStock,
    }

    return response, nil
}

func (h *CartGrpcHandler) ClearCart(ctx context.Context,req *pb.ClearCartRequest)(*pb.ClearCartResponse,error){

    err := h.cartService.ClearCart(ctx,req.UserID)

    if err != nil {
        return nil,err
    }

    return &pb.ClearCartResponse{Done: true},nil
}