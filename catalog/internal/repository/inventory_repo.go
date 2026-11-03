package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
)

type InventoryRepo struct {
	DB *pgxpool.Pool
}

func InitInventoryRepo(db *pgxpool.Pool) *InventoryRepo{
	return &InventoryRepo{DB: db}
}

func (repo *InventoryRepo) ReserveInvetory(ctx context.Context,req models.ReservationRequest)(bool,error){


	tx,err := repo.DB.Begin(ctx)

	if err != nil {
		return false,nil
	}

	defer tx.Rollback(ctx)

	query := `UPDATE products SET stock_quantity = stock_quantity - $1 WHERE id = $2 AND stock_quantity >= $1 RETURNING id;`

	for _,item := range req.Items {
		
		if item.Quantity <= 0 {
			return false,fmt.Errorf("quantity should greater than 0")
		}

		var productId int64

		err := tx.QueryRow(ctx,query,item.Quantity,item.ProductID).Scan(&productId)

		if err != nil {
            if errors.Is(err, pgx.ErrNoRows) {
                return false, fmt.Errorf(
                    "product %d does not exist or has insufficient stock",
                    item.ProductID,
                )
            }

			return false,err
		}
	}

	var rows[][]interface{}

	for _,p := range req.Items {
		rows = append(rows, []interface{}{req.OrderID,p.ProductID,p.Quantity})
	}

	_,err = repo.DB.CopyFrom(ctx,pgx.Identifier{"inventory_reservations"},[]string{"order_id","product_id","quantity"},pgx.CopyFromRows(rows))

	if  err != nil {
		return false,err
	}

	if err = tx.Commit(ctx); err != nil {
		return false,err
	}

	return true,nil
}

func (repo *InventoryRepo)ReleaseInventory(ctx context.Context,orderID int64)(error){

	tx,err := repo.DB.Begin(ctx)
	if err != nil {

		return err
	}
	defer tx.Rollback(ctx)

	query := `SELECT product_id,quantity FROM inventory_reservations WHERE order_id = $1 AND status = 'RESERVED' FOR UPDATE;`

	rows,err := tx.Query(ctx,query,orderID)

	if err != nil {
		return err
	}

	var products []models.ReservationItem

	for rows.Next(){

		var product models.ReservationItem

		err := rows.Scan(&product.ProductID,&product.Quantity)

		if err != nil {
			return err
		}

		products = append(products, product)
	}

	defer rows.Close()


    if len(products) == 0 {
        return  errors.New("no active reservation found")
    }

	updateQuery := `UPDATE products SET stock_quantity = stock_quantity + $1 WHERE id = $2;`

	for _,p := range products {

		_,err := tx.Exec(ctx,updateQuery,p.Quantity,p.ProductID)

		if err != nil {
			return err
		}


	}

	updateInventoryQuery := `UPDATE inventory_reservations SET status = 'RELEASED' WHERE order_id = $1 AND status = 'RESERVED'`

	_,err = tx.Exec(ctx,updateInventoryQuery,orderID)

	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil

}


func (repo *InventoryRepo) ConfirmInventory(ctx context.Context,orderID int64)(error){
	query := `UPDATE inventory_reservations SET status = 'CONFIRMED WHERE order_id = $1 AND status = 'RESERVED` 

	tx,err := repo.DB.Begin(ctx)

	if err != nil {
		return err
	}

	_,err = tx.Exec(ctx,query,orderID)

	if err != nil {
		return err
	}

	tx.Rollback(ctx)

	if err = tx.Commit(ctx); err != nil{
		return err
	}

	return nil

}