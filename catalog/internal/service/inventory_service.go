package service

import (
	"context"
	"time"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/repository"
)

type InventoryService struct {
	repo *repository.InventoryRepo
}

func InitInventoryService(repo *repository.InventoryRepo) *InventoryService{
	return &InventoryService{repo: repo}
}

func (s *InventoryService) ReserveInventory(ctx context.Context,req models.ReservationRequest)(error){

	if len(req.Items) == 0 {
		return ErrInvalidProduct
	}

	expiresAt := time.Now().Add(10 * time.Minute)

	_,err := s.repo.ReserveInvetory(ctx,req,expiresAt)

	if err != nil{
		return err
	}

	return nil
}