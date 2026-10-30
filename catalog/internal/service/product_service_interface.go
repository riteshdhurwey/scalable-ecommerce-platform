package service

import (
	"context"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/dto"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
)


type ProductInterface interface {
    CreateProduct(ctx context.Context, input dto.ProductRequestDTO) (*models.Product, error)

    GetProductByID(ctx context.Context, productID int64) (*models.Product, error)

    ListProducts(ctx context.Context, page int, limit int) ([]*models.Product, error)

    ListProductsByCategory(ctx context.Context, categoryID int64, page int, limit int) (*[]models.Product, error)

    UpdateProduct(ctx context.Context, input dto.UpdateProductDTO) (*models.Product, error)

    DeleteProduct(ctx context.Context, productID int64) error

    UpdateInventory(ctx context.Context, productID int64, quantity int) error

    ValidateStock(ctx context.Context, productID int64, requestedQuantity int) error
}
