package service

import (
	"context"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/repository"
)

type CategoryService struct {
	repo repository.PostgreDB
}

func InitCategoryService(repo repository.PostgreDB) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, categoryName string, categoryDesc string) (int64, error) {

	category := models.Category{Name: categoryName, Description: categoryDesc}

	categoryID, err := s.repo.CreateCategory(ctx, category)

	if err != nil {
		return 0, err
	}

	return categoryID, nil
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]models.Category, error) {

	cotegories, err := s.repo.GetCategories(ctx)

	if err != nil {
		return nil, err
	}

	return cotegories, nil
}
