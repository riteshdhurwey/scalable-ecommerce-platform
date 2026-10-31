package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/repository"
)

var (

	ErrCartIDNotFound = errors.New("Cart Id not found")
	ErrNoItemFound    = errors.New("No item found")
)

type CartSevice struct {
	repo *repository.CartRepo
}

func InitCartSevice(repo *repository.CartRepo) *CartSevice {
	return &CartSevice{repo: repo}
}

func (s *CartSevice) AddToCart(ctx context.Context,userID int64,req models.Cart)(error){

	if req.ProductID <= 0 || req.Quantity <= 0 {
		return fmt.Errorf("invalid product")
	}

	productID  := strconv.Itoa(int(req.ProductID))
	catalogURL := os.Getenv("CATALOG_SERVICE_URL")
    url := catalogURL + "/api/v1/product/" + productID

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

	if product.StockQuantity < req.Quantity {
		return fmt.Errorf("available stock : %d",product.StockQuantity)
	}

	cartID,err := s.repo.IsCartExists(ctx,userID)

	if cartID == nil {

		newCartID , err := s.repo.CreateCart(ctx,userID)

		if err != nil {
			return err
		}

		cartID = &newCartID
	}

	err = s.repo.AddToCart(ctx,*cartID,req)

	if err != nil {
		return err
	}

	
	return nil

}


func (s *CartSevice) GetCart(ctx context.Context,userID int64)(*[]models.Product,error){

	cart_id,err := s.repo.IsCartExists(ctx,userID)

	if err != nil {
		return nil,err
	}

	if cart_id == nil {
		return nil,ErrCartIDNotFound
	}

	cartItems,err := s.repo.GetCart(ctx,*cart_id)

	if err != nil {
		return nil,err
	}

	if cartItems == nil {
		return nil,ErrNoItemFound
	}
	
	var productIDs []int64
	productIDMap := make(map[int64]int)
	
	for _,product :=  range *cartItems {
		
		productIDs = append(productIDs, product.ProductID)
		productIDMap[product.ProductID] = product.Quantity
	}
	
	body := map[string][]int64{
		"product_ids":productIDs,
	}
	
	
	jsonData , err := json.Marshal(&body)
	
	if err != nil {
		return nil,err
	}
	
    catalogURL := os.Getenv("CATALOG_SERVICE_URL")
	
	url := catalogURL + "/api/v1/products/batch"
	
	request,err := http.NewRequest(http.MethodPost,url,bytes.NewBuffer(jsonData))
	
	if err != nil {
		return nil,err
	}
	
	var httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp,err := httpClient.Do(request)

	if err != nil {
		return nil,err
	}
	
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"catalog returned status %d",
			resp.StatusCode,
		)
	}
	
	var response []models.BatchRequest

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil,err
	}
	var products []models.Product
	
	for _,value := range response {

		var product models.Product

		product.ProductID = value.ID
		product.Name  = value.Name
		product.Price = value.Price
		product.StockQuantity = productIDMap[value.ID]

		products = append(products, product)
	}


	return &products,nil


}