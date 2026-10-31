package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/models"
)

type CartRepo struct {
	DB *pgxpool.Pool
}

func InitCartRepo(db *pgxpool.Pool) *CartRepo {
	return &CartRepo{DB: db}
}

func (repo *CartRepo) CreateCart(ctx context.Context, user_id int64) (int64, error) {

	query := `INSERT INTO carts (user_id) VALUES($1) RETURNING cart_id;`
	var cartID int64
	err := repo.DB.QueryRow(ctx, query, user_id).Scan(&cartID)

	if err != nil {
		return 0, err
	}

	return cartID, nil

}

func (repo *CartRepo) IsCartExists(ctx context.Context, user_id int64) (*int64, error) {
	query := `SELECT cart_id FROM carts WHERE user_id = $1;`

	var cartID int64

	err := repo.DB.QueryRow(ctx, query, user_id).Scan(&cartID)

	if err != nil {
		return nil, err
	}

	return &cartID, nil
}

func (repo *CartRepo) AddToCart(ctx context.Context, cartID int64, req models.Cart) error {
	query := `INSERT INTO cart_items (cart_id, product_id, quantity)
				  VALUES ($1, $2, $3)
				  ON CONFLICT (cart_id, product_id)
				  DO UPDATE
				  SET quantity = cart_items.quantity + EXCLUDED.quantity;`

	err := repo.DB.QueryRow(ctx,query,cartID,req.ProductID,req.Quantity)

	if err != nil {
		return ctx.Err()
	}

	return nil
}


func (repo *CartRepo) GetCart(ctx context.Context,cartID int64)(*[]models.Cart,error){
	
	query := `SELECT product_id,quantity FROM cart_items WHERE cart_id = $1;`

	row,err := repo.DB.Query(ctx,query,cartID)

	if err != nil {
		return nil,err
	}

	var products []models.Cart

	for row.Next(){

		var product models.Cart 

		err := row.Scan(&product.ProductID,&product.Quantity)

		if err != nil {
			return nil,err
		}

		products = append(products, product)
	}

	defer row.Close()

	return &products,nil
}


func (repo *CartRepo) UpdateProduct(ctx context.Context,cartID int64,req models.Cart)(bool,error){

	query := `UPDATE cart_items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3;`

	result,err := repo.DB.Exec(ctx,query,req.Quantity,cartID,req.ProductID)

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		return false,nil
	}

	return true,nil
}


func (repo *CartRepo) DeleteItem(ctx context.Context,cartID int64,productID int64)(bool,error){

	query := `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2;`

	result,err := repo.DB.Exec(ctx,query,cartID,productID)

	if err != nil {
		return false,err
	}

	if result.RowsAffected() == 0 {
		return false,nil
	}

	return true,nil
}