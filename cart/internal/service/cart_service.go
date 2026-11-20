package service

import (

	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	catalogclient "github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/clients/catalog"
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
	grpcClient catalogclient.Client
}

func InitCartSevice(repo *repository.CartRepo,grpcClient catalogclient.Client) *CartSevice {
	return &CartSevice{repo: repo,grpcClient: grpcClient}
}

func (s *CartSevice) AddToCart(ctx context.Context,userID int64,req models.Cart)(*models.Product,error){

	if req.ProductID <= 0 || req.Quantity <= 0 {
		return nil, fmt.Errorf("invalid product")
	}

	var productIDs []int64

	productIDs = append(productIDs, req.ProductID)

	product,err := s.grpcClient.GetProduct(ctx,req.ProductID)

	if err != nil {
		log.Printf("Catalog gRPC error: %v", err)
		return nil,err
	}
	
	cartID, err := s.repo.IsCartExists(ctx, userID)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("checking cart: %w", err)
	}

	currentQuantity := 0

	if cartID != nil {
		currentQuantity, err = s.repo.GetCartProductByID(
			ctx, *cartID, req.ProductID,
		)

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("checking cart item: %w", err)
		}
	}

	newQuantity := req.Quantity + currentQuantity

	if int(product.Quantity) < newQuantity {
		return nil,fmt.Errorf("available stock : %d and Ordered Quantity is %d",product.Quantity,newQuantity)
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

	productResponse := models.Product{Name: product.Name,ProductID: req.ProductID,Price: product.Price,StockQuantity: req.Quantity}

	
	return &productResponse, nil

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


	product,err := s.grpcClient.GetProduct(ctx,req.ProductID)

	if err != nil {
		return false,err
	}

	newQuantity := previousQuantity + req.Quantity

	if newQuantity > int(product.Quantity) {
		return false,fmt.Errorf("available stock : %d",product.Quantity)

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


func (s *CartSevice) CheckOut(ctx context.Context,userID int64)(float32,*[]models.Product,[]models.OutOfStockProducts,error){

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

	log.Println("got cart items")
	if err != nil {
		return 0,nil,nil,err
	}

	
	var productIDs []int64
	productIdOrderedStock := make(map[int64]int)
	
	for _,product :=  range *cartItems {
		
		productIDs = append(productIDs, product.ProductID)
		productIdOrderedStock[product.ProductID] = product.Quantity
	}
	
	response ,err := s.grpcClient.GetProductByIDs(ctx,productIDs)

	if err != nil {
		return 0,nil,nil,err
	}

	var products []models.Product
	var outOFStockProducts []models.OutOfStockProducts
	var totalAmount float32

	for _,p := range response.Products {
		

		if productIdOrderedStock[p.Id] > int(p.Quantity) {
			var product models.OutOfStockProducts
			
			product.ProductID = p.Id
			product.OrderedStock = productIdOrderedStock[p.Id]
			product.Name = p.Name
			product.AvailableStock = int(p.Quantity)
			
			outOFStockProducts = append(outOFStockProducts, product)
			continue
		}

		var product models.Product


			product.ProductID = p.Id
			product.Name = p.Name
			product.Price = p.Price
			product.StockQuantity = productIdOrderedStock[p.Id]

			totalAmount += float32(p.Price) * float32(product.StockQuantity)

			products = append(products, product)

	}

	return totalAmount,&products,outOFStockProducts,nil


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