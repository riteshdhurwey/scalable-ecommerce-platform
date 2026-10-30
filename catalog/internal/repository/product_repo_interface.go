package repository

import (
	"context"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/dto"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
)

type ProductRepository interface {
    CreateProduct(ctx context.Context, product *dto.ProductRequestDTO) error
    GetProductByID(ctx context.Context, productID int64) (*models.Product, error)
    ProductList(ctx context.Context, limit, offset int) (*[]models.Product, error)
    ListByCategory(ctx context.Context, categoryID int64, limit, offset int) (*[]models.Product, error)
    Update(ctx context.Context, product *dto.UpdateProductDTO) error
    DeleteProduct(ctx context.Context, productID int64) error
   // UpdateStock(ctx context.Context, productID int64, quantity int) error
	CategoryExists(ctx context.Context,categoryID int64) (*int64,string,error)
}