package models

type Cart struct {
	ProductID int64 `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Product struct {
    ProductID     int64  `json:"product_id"`
    Name          string `json:"name"`
    Price         float64  `json:"price"`
    StockQuantity int  `json:"stock_quantity"`
}


type BatchRequest struct {
    ID  int64 `json:"id"`
    Name string `json:"name"`
    Price float64 `json:"price"`
}