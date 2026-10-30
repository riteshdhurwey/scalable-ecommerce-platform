package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/dto"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/repository"
)


var (
    ErrProductNotFound  = errors.New("product not found")
    ErrInvalidProduct   = errors.New("invalid product")
    ErrInvalidPrice     = errors.New("invalid price")
    ErrInvalidStock     = errors.New("invalid stock quantity")
    ErrInvalidCategory  = errors.New("invalid category")
    ErrCategoryNotFound = errors.New("category not found")
	ErrInvalidInput 	= errors.New("Invalid input")
	ErrInvalidPageSize  = errors.New("Invalid page or pageSize")
)

type ProductService struct {
	repo * repository.PostgreDB
}

func InitProductService(repo *repository.PostgreDB) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context,req dto.ProductRequestDTO) (*models.Product,error){

	if req.Price <= 0 {
		return nil,ErrInvalidPrice
	}

	if req.StockQuantity < 0 {
		return nil,ErrInvalidStock
	}

	catoryName,err := s.repo.CategoryExists(ctx,req.CategoryID)

	if err != nil {
		return nil,fmt.Errorf("check category existence: %w", err)
	}

	if catoryName == nil {
		return nil,ErrCategoryNotFound
	}


	product,err := s.repo.CreateProduct(ctx,req)

	if err != nil {
		return nil,err
	}

	return &product,nil
	
}


func (s *ProductService) GetProductByID (ctx context.Context,productID int64)(*models.Product,error){

	product,err := s.repo.GetProductByID(ctx,productID)

	if err != nil {
		return nil,err
	}

	return &product,nil
}


func (s *ProductService) DeleteProduct(ctx context.Context,productID int64) (bool,error){

	delete,err := s.repo.DeleteProduct(ctx,productID)

	if err != nil {
		return false,err
	}

	if !delete {
		return false,ErrProductNotFound
	}

	return true,nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, input dto.UpdateProductDTO) (*models.Product, error){
	if input.ID == nil || (input.Price == nil && input.StockQuantity == nil){
		return nil,ErrInvalidInput
	}

	err := s.repo.Update(ctx,&input)

	if err != nil {
		return nil,err
	}

	product,err := s.repo.GetProductByID(ctx,*input.ID)

	if err != nil {
		return nil,fmt.Errorf("product updated failed to get product details %w",err)
	}

	return &product,nil

} 

func(s *ProductService) ListProducts(ctx context.Context, page int, limit int) (*[]models.Product, error){
		maxLimit := 50

		if page < 0 || limit > maxLimit {
			return nil,ErrInvalidPageSize
		}

		offset := (page -1 )*limit

		products,err := s.repo.ProductList(ctx,limit,offset)

		if err != nil {
			return nil,err
		}

		return products,nil
}


func (s *ProductService) ListProductsByCategory(ctx context.Context,categoryID int64,page int,limit int)(*[]models.Product,string,error){

	maxLimit := 50

	if page < 0 || limit > maxLimit {
		return nil,"",ErrInvalidPageSize
	}

	offset := (page - 1)*limit

	categoryName,err := s.repo.CategoryExists(ctx,categoryID)

	if err != nil {
		return nil,"",fmt.Errorf("check category existence: %w", err)
	}

	if categoryName == nil {
		return nil,"",ErrCategoryNotFound
	}

	products,err := s.repo.ListByCategory(ctx,categoryID,limit,offset)

	if err != nil {
		return nil,"",err
	}

	if products == nil {
		return nil,"",ErrProductNotFound
	}

	return products,*categoryName,err
}





