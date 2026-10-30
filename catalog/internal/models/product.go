package models

import "time"


type Product struct {
	ID 			  int64  `json:"id"`
    Name          string `json:"name"`
    Description   string `json:"description"`
    Price         float64 `json:"price"`
    StockQuantity int32  `json:"stock_quantity"`
	Created_At    time.Time `json:"created_at"`
}