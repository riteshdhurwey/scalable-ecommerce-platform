package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/models"
)


type OrderRepo struct {
	DB *pgxpool.Pool
}

func InitOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{DB: db}
}

func (repo *OrderRepo) CreateOrder(ctx context.Context,req models.Order)(*int64,error){
	query := `INSERT INTO orders (user_id,total_amount,payment_expires_at) VALUES($1,$2,$3) RETURNING id;`

	var OrderID int64

	err := repo.DB.QueryRow(ctx,query,req.UserID,req.TotalAmount,req.PaymentExpAt).Scan(&OrderID)

	if err != nil {
		return nil,err
	}

	return &OrderID,nil

}


func(repo *OrderRepo) CreateOrderItems(ctx context.Context,OrderID int64,products []models.OrderItems)(error){

	var rows [][]interface{}

	for _,p:= range products {

		rows = append(rows, []interface{}{OrderID,p.ProductID,p.ProductName,p.Quantity,p.UnitPrice,p.SubTotal})
	}

	_, err := repo.DB.CopyFrom(ctx,pgx.Identifier{"order_items"},[]string{"order_id","product_id","product_name","quantity","unit_price","subtotal"},pgx.CopyFromRows(rows),)

	return err
}






