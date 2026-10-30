package dto


type UpdateProductDTO struct {
	ID *int64 `josn:"id"`
	Price *int64 `json:"price"`
	StockQuantity *int `json:"stock_quantity"` 
}

type ProductRequestDTO struct {
    Name          string `json:"name" binding:"required,max=255"`
    Description   string `json:"description"`
    Price         float64 `json:"price" binding:"required,gt=0"`
    StockQuantity int32  `json:"stock_quantity" binding:"gte=0"`
    CategoryID    int64  `json:"category_id" binding:"required,gt=0"`
}
