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

	query := `INSERT INTO carts (user_id) VALUES($1) RETURINING cart_id;`
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

func (repo *CartRepo) AddToCart(ctx context.Context, cartID int64, req models.CART) error {
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

