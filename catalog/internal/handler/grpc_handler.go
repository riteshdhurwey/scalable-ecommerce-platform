package handler

import (
	"context"

	pb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/catalog"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogHandler struct {
    pb.UnimplementedCatalogServiceServer

    catalogService   *service.ProductService
    inventoryService *service.InventoryService
}

func InitCatalogGRPCHandler(catalogservice *service.ProductService,inventoryservice *service.InventoryService) *CatalogHandler {
	return &CatalogHandler{catalogService: catalogservice,inventoryService: inventoryservice}
}

func (h *CatalogHandler) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {

    products, err := h.catalogService.GetProductsByIDs(
        ctx,
        req.Id,
    )
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to retrieve products")
    }

    response := &pb.GetProductsResponse{
        Products: make([]*pb.Product, 0, len(*products)),
    }

    for _, product := range *products {
        response.Products = append(response.Products, &pb.Product{
            Id:            product.ProductID,
            Name:          product.Name,
            Price:         product.Price,
            Quantity: 	   int32(product.AvailableStock),
        })
    }

    return response, nil


}

func (h *CatalogHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest)(*pb.Product,error){

	if req.Id <= 0 {
		return nil,status.Error(codes.InvalidArgument,"invalid product id")
	}

	product,err := h.catalogService.GetProductByID(ctx,req.Id)

	if err != nil {
		return nil,status.Error(codes.NotFound,"product not found")
	}

	return &pb.Product{Id: product.ID,Name: product.Name,Price: product.Price,Quantity: product.StockQuantity},nil
}
