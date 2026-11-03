package models
type ReservationRequest struct {
    OrderID int64
    Items   []ReservationItem
}

type ReservationItem struct {
    ProductID int64   `json:"product_id"`
    Quantity  int     `json:"quantity"`
}
