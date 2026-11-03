package models

import "time"

type ReservationRequest struct {
    OrderID int64
    Items   []ReservationItem
}

type ReservationItem struct {
    ProductID int64   `json:"product_id"`
    Quantity  int     `json:"quantity"`
}


type Product struct {
    ProductID     int64  `json:"product_id"`
    Name          string `json:"name"`
    Price         float64  `json:"price"`
    OrderQuantity int  `json:"stock_quantity"`
}

type Order struct {
	UserID int64 `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	PaymentExpAt time.Time `json:"payment_expires_at"`
}

type OrderItems struct {
	ProductID int64 `json:"product_id"`
	ProductName string `json:"name"`
	Quantity int  `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	SubTotal float64 `json:"subtotal"`
}


type OrderResponse struct {
	UserID int64 `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	PaymentStatus string `json:"payment_status"`
	PaymentExpAt time.Time `json:"payment_expires_at"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type OutOfStockProducts struct {
	ProductID int64 `json:"product_id"`
    Name string `json:"name"`
	OrderedStock int `json:"ordered_stock"`
    AvailableStock int `json:"available_stock"`
}

type OrderItemResponse struct {
	OrderID  int64 	`json:"order_id"`
	ProductID int64 `json:"product_id"`
	ProductName string `json:"name"`
	Quantity int  `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	SubTotal float64 `json:"subtotal"`
	Created_At time.Time `json:"created_at"`
}


type ProductPayload struct {
	TotalAmount        float64   `json:"totalAmount"`
	AvailableProducts  []Product `json:"available"`
	OutOfStockProducts []OutOfStockProducts `json:"out of stock"` 
}