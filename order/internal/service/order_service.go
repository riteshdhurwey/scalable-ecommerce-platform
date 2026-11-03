package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/repository"
)

var (

	ErrOutOfStock = errors.New("some products are out of stock")
	ErrEmptyCart  = errors.New("Cart is empty")
)


type OrderService struct {
	repo *repository.OrderRepo
}

func InitOrderService(repo *repository.OrderRepo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CheckOut(ctx context.Context,userID int64,token string)(*[]models.OutOfStockProducts,*[]models.Product,float32,error){

	cartURL := os.Getenv("CART_SERVICE_URL")

	url := cartURL + "/api/v1/cart/checkout"

	request,err := http.NewRequestWithContext(ctx,http.MethodGet,url,nil)

	request.AddCookie(&http.Cookie{
		Name:  "authorization",
		Value: token,
	})

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil,nil,0,err
	}

	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusConflict{
			return nil,nil,0,ErrEmptyCart
		}
		return nil, nil, 0, fmt.Errorf(
			"cart service returned status %d:",
			resp.StatusCode,
		)
	}
	var  response models.ProductPayload

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil,nil,0,err

	}


	if response.OutOfStockProducts != nil {
		return &response.OutOfStockProducts,nil,0,ErrOutOfStock
	}

	if response.AvailableProducts == nil {
		return nil,nil,0,ErrEmptyCart
	}

	order := models.Order{UserID: userID,TotalAmount: response.TotalAmount,PaymentExpAt: time.Now().Add(10 * time.Minute)}

	orderID,err := s.repo.CreateOrder(ctx,order)

	if err != nil {
		return nil,nil,0,err
	}

	if orderID == nil {
		return nil,nil,0,fmt.Errorf("failed to create order")
	}

	var orderReq models.ReservationRequest
	orderReq.OrderID = *orderID

	var items []models.ReservationItem
	
	AvailableProducts := response.AvailableProducts

	for _,p := range AvailableProducts {
		var item models.ReservationItem

		item.ProductID = p.ProductID
		item.Quantity = p.OrderQuantity

		items = append(items, item)
	}

	orderReq.Items = items

	jsonData,err := json.Marshal(&orderReq)

	if err != nil{
		return nil,nil,0,err
	}

	 catalogBaseURL := os.Getenv("CATALOG_SERVICE_URL")

	catalogURL := catalogBaseURL + "/api/v1/inventory/reserve"

	reserveReq,err := http.NewRequest(http.MethodPost,catalogURL,bytes.NewBuffer(jsonData))

	if err != nil{

		return nil,nil,0,err
	}

	reserveResp,err := http.DefaultClient.Do(reserveReq)

	if err != nil{
		return nil,nil,0,err
	}

	defer reserveResp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusConflict{
			return nil,nil,0,ErrEmptyCart
		}
		return nil, nil, 0, fmt.Errorf(
			"cart service returned status %d: ",
			resp.StatusCode,
		)
	}
	var products[]models.OrderItems

	for _,p := range AvailableProducts {

		var product models.OrderItems

		product.ProductID = p.ProductID
		product.ProductName = p.Name
		product.Quantity = p.OrderQuantity
		product.UnitPrice = p.Price
		product.SubTotal = p.Price * float64(p.OrderQuantity)
		products = append(products, product)

	}

	err = s.repo.CreateOrderItems(ctx,*orderID,products)

	if err != nil {
		return nil,nil,0,err
	}

	clearCartURL := cartURL + "/api/v1/cart/delete"

	clearReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		clearCartURL,
		nil,
	)
	if err != nil {
		return nil, nil, 0, err
	}

	clearReq.AddCookie(&http.Cookie{
		Name:  "authorization",
		Value: token,
	})

	clearResp, err := http.DefaultClient.Do(clearReq)
	if err != nil {
		return nil,nil,0,fmt.Errorf("cart not cleared")
	}

	defer clearResp.Body.Close()

	if clearResp.StatusCode != http.StatusOK {
			return nil,nil,0,fmt.Errorf("cart returned status code %d",clearResp.StatusCode)
	}


	return nil,&AvailableProducts,float32(order.TotalAmount),nil

	
	
}