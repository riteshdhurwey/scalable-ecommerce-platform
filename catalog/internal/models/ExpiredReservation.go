package models

type ExpiredReservation struct {
	ID      int64 `json:"id"`
	OrderID int64 `json:"order_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int `json:"quantity"`
}