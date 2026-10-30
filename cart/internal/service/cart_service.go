package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/repository"
)

type CartSevice struct {
	repo *repository.CartRepo
}

func InitCartSevice(repo *repository.CartRepo) *CartSevice {
	return &CartSevice{repo: repo}
}

func (s *CartSevice) AddToCart(ctx context.Context,userID int64,req models.CART)(error){

	var catalogBaseURL = os.Getenv("CATALOG_SERVICE_URL")
	url := fmt.Sprintf("%sproduct/%d",catalogBaseURL,req.ProductID)

	request,err := http.NewRequestWithContext(ctx,http.MethodGet,url,nil)

	if err != nil {
		return err
	}
	var httpClient = &http.Client{
   		 Timeout: 3 * time.Second,
	}


	resp,err := httpClient.Do(request)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return  fmt.Errorf(
            "catalog returned status %d",
            resp.StatusCode,
        )
	}

	var product models.Product

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return err
	}

	if product.StockQuantity < int64(req.Quantity) {
		return fmt.Errorf("available stock : %d",product.StockQuantity)
	}

	cartID,err := s.repo.IsCartExists(ctx,userID)

	if cartID == nil {
		*cartID , err = s.repo.CreateCart(ctx,userID)

		if err != nil {
			return err
		}
	}

	err = s.repo.AddToCart(ctx,*cartID,req)

	if err != nil {
		return err
	}

	


}
