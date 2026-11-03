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
	ErrEmptyCart    = errors.New("cart is emtpy")
	ErrProductNotFound = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient Stock")
)

type CartSevice struct {
	repo *repository.CartRepo
}

func InitCartSevice(repo *repository.CartRepo) *CartSevice {
	return &CartSevice{repo: repo}
}

func (s *CartSevice) AddToCart(ctx context.Context,userID int64,req models.Cart)(*models.Product,error){

	if req.ProductID <= 0 || req.Quantity <= 0 {
		return nil, fmt.Errorf("invalid product")
	}

	productID  := strconv.Itoa(int(req.ProductID))
	catalogURL := os.Getenv("CATALOG_SERVICE_URL")
    url := catalogURL + "/api/v1/product/" + productID

	request,err := http.NewRequestWithContext(ctx,http.MethodGet,url,nil)

	if err != nil {
		return nil,err
	}
	var httpClient = &http.Client{
   		 Timeout: 3 * time.Second,
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

	var product models.Product

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil,err
	}

	
	cartID,err := s.repo.IsCartExists(ctx,userID)
	
	currentQuantity := 0
	
	if cartID != nil {
		
		currentQuantity,_ = s.repo.GetCartProductByID(ctx,*cartID,req.ProductID)
		
	}

	newQuantity := req.Quantity + currentQuantity

	if product.StockQuantity < newQuantity {
		return nil,fmt.Errorf("available stock : %d",product.StockQuantity)
	}

	if cartID == nil {

		newCartID , err := s.repo.CreateCart(ctx,userID)

		if err != nil {
			return nil,err
		}

		cartID = &newCartID
	}


	

	err = s.repo.AddToCart(ctx,*cartID,req)

	if err != nil {
		return nil,err
	}

	productResponse := models.Product{Name: product.Name,ProductID: product.ProductID,Price: product.Price,StockQuantity: req.Quantity}

	
	return &productResponse, nil

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
		return nil,ErrEmptyCart
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


func (s *CartSevice) UpdateCart(ctx context.Context,userID int64,req models.Cart)(bool,error){

	if req.Quantity <= 0 {
		return false,fmt.Errorf("invalid Quantity")
	}

	cartId,err := s.repo.IsCartExists(ctx,userID)

    if err != nil {
		return false,err
	}

	if cartId == nil {
		return false,ErrCartIDNotFound
	}

	previousQuantity,err := s.repo.GetCartProductByID(ctx,*cartId,req.ProductID)

	if err != nil{
		return false,err
	}

	productID  := strconv.Itoa(int(req.ProductID))
	catalogURL := os.Getenv("CATALOG_SERVICE_URL")
    url := catalogURL + "/api/v1/product/" + productID

	request,err := http.NewRequestWithContext(ctx,http.MethodGet,url,nil)

	if err != nil {
		return false,err
	}
	var httpClient = &http.Client{
   		 Timeout: 3 * time.Second,
	}


	resp,err := httpClient.Do(request)

	if err != nil {
		return false,err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf(
            "catalog returned status %d",
            resp.StatusCode,
        )
	}

	var product models.Product

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return false,err
	}

	newQuantity := previousQuantity + req.Quantity

	if newQuantity > product.StockQuantity {
		return false,fmt.Errorf("available stock : %d",product.StockQuantity)

	}

	updated,err := s.repo.UpdateProduct(ctx,*cartId,req)

	if err != nil {
		return false,err
	}

	if !updated {
		return false,ErrProductNotFound

	}

	return true,nil
}

func(s *CartSevice) DeleteItem(ctx context.Context,userID,productID int64)(error){

	cartID,err := s.repo.IsCartExists(ctx,userID)

	if err != nil {
		return err
	}

	if cartID == nil {

		return ErrCartIDNotFound
	}

	deleted,err := s.repo.DeleteItem(ctx,*cartID,productID)

	if err != nil {
		return err
	}

	if !deleted {
		return ErrProductNotFound
	}

	return nil
}


func (s *CartSevice) CheckOut(ctx context.Context,userID int64)(float32,*[]models.Product,*[]models.OutOfStockProducts,error){

	cart_id,err := s.repo.IsCartExists(ctx,userID)

	if cart_id == nil {
		return 0,nil,nil,ErrCartIDNotFound
	}

	if err != nil {
		return 0,nil,nil,err
	}

	cartItems,err := s.repo.GetCart(ctx,*cart_id)

	if cartItems == nil{
		return 0,nil, nil,ErrEmptyCart
	}
	if err != nil {
		return 0,nil,nil,err
	}

	
	var productIDs []int64
	productIdOrderedStock := make(map[int64]int)
	
	for _,product :=  range *cartItems {
		
		productIDs = append(productIDs, product.ProductID)
		productIdOrderedStock[product.ProductID] = product.Quantity
	}
	
	body := map[string][]int64{
		"product_ids":productIDs,
	}
	
	
	jsonData , err := json.Marshal(&body)
	
	if err != nil {
		return 0,nil,nil,err
	}
	
    catalogURL := os.Getenv("CATALOG_SERVICE_URL")
	
	url := catalogURL + "/api/v1/products/batch"
	
	request,err := http.NewRequest(http.MethodPost,url,bytes.NewBuffer(jsonData))
	
	if err != nil {
		return 0,nil,nil,err
	}
	
	var httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp,err := httpClient.Do(request)

	if err != nil {
		return 0,nil, nil,err
	}
	
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0,nil, nil,fmt.Errorf(
			"catalog returned status %d",
			resp.StatusCode,
		)
	}
	
	var response []models.BatchRequest

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0,nil,nil,err
	}
	var products []models.Product
	var outOFStockProducts []models.OutOfStockProducts
	var totalAmount float32

	for _,p := range response {
		
		if productIdOrderedStock[p.ID] > p.AvailableStock {
			var product models.OutOfStockProducts
			
			product.ProductID = p.ID
			product.OrderedStock = productIdOrderedStock[p.ID]
			product.Name = p.Name
			product.AvailableStock = p.AvailableStock
			
			outOFStockProducts = append(outOFStockProducts, product)
			continue
		}

		var product models.Product

			product.ProductID = p.ID
			product.Name = p.Name
			product.Price = p.Price
			product.StockQuantity = productIdOrderedStock[p.ID]

			totalAmount += float32(p.Price) * float32(product.StockQuantity)

			products = append(products, product)

	}

	if products == nil {
		return 0,nil,nil,ErrEmptyCart
	}

	return totalAmount,&products,&outOFStockProducts,nil


}

func (s *CartSevice) ClearCart(ctx context.Context, userID int64) error {
    cartID, err := s.repo.IsCartExists(ctx, userID)
    if err != nil {
        return err
    }

    if cartID == nil {
        return ErrCartIDNotFound
    }

    return s.repo.ClearCart(ctx, *cartID)
}