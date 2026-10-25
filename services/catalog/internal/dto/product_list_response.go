package dto

import "github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"

type ProductListResponse struct{
	Page int `json:"page"`
	Limit int `json:"limit"`
	CategoryID int64 `json:"category_id"`
	CategoryName string `json:"category_name"`
	Products *[]models.Product
}