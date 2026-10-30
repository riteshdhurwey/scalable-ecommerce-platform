package models

type CART struct {
	ProductID int64 `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Product struct {
    ProductID     int64  `json:"product_id"`
    Name          string `json:"name"`
    Price         int64  `json:"price"`
    StockQuantity int64  `json:"stock_quantity"`
}