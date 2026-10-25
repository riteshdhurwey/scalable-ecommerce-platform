package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/dto"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
)

var (
	ErrProductNotFound  = errors.New("product not found")
)


func (repo *PostgreDB) CreateProduct(ctx context.Context,p dto.ProductRequestDTO)(models.Product,error){
	query := `INSERT INTO products (name,description,price,stock_quantity,category_id)
			   VALUES($1,$2,$3,$4,$5) RETURNING id,created_at;`
	
	var productID int64
	var created_At time.Time

	err := repo.DB.QueryRow(ctx,query,p.Name,p.Description,p.Price,p.StockQuantity,p.CategoryID).Scan(&productID,&created_At)

	if err != nil {
		return models.Product{},err
	}

	product := models.Product{ID: productID,Name: p.Name,Description: p.Description,Price: p.Price,StockQuantity: p.StockQuantity,Created_At: created_At}

	return product,nil
}


func (repo *PostgreDB) CategoryExists(ctx context.Context,categoryID int64)(*string,error){
	query := `SELECT name FROM categories WHERE id = $1;`

	var name string
	err := repo.DB.QueryRow(ctx,query,categoryID).Scan(&name)

	if err != nil {
		return nil,err
	}

	return &name,nil
}


func (repo *PostgreDB) GetProductByID(ctx context.Context,productID int64) (models.Product,error){
	query := `SELECT id,name,description,price,stock_quantity,created_at FROM products WHERE id = $1;`

	var product models.Product

	err := repo.DB.QueryRow(ctx,query,productID).Scan(&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.Created_At,
		)	
		
	if err != nil {
		return models.Product{},err
	}

	return product,nil
}


func (repo *PostgreDB) DeleteProduct(ctx context.Context,productID int64) (bool,error){
	query := `DELETE FROM Products WHERE id = $1;`

	result,err := repo.DB.Exec(ctx,query)

	if err != nil {
		return false,err
	}

	if result.RowsAffected() <= 0 {
		return false,nil
	}

	return true,nil
}


func (repo * PostgreDB) Update(ctx context.Context,product *dto.UpdateProductDTO) error {
	var sets[]string
	var args []any

	if product.Price != nil {
		args = append(args, product.Price)
		sets = append(sets, fmt.Sprintf("price = $%d",len(args)))
	}

	if product.StockQuantity != nil {
		args = append(args, product.StockQuantity)
		sets = append(sets, fmt.Sprintf("stock_quantity = $%d",len(args)))
	}

	if len(sets) == 0 {
		return fmt.Errorf("No field to Update")
	}

	args = append(args, product.ID)

	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d",strings.Join(sets,", "),len(args),)

	result,err := repo.DB.Exec(ctx,query,args...)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrProductNotFound
	}

	return nil
}

func (repo *PostgreDB) ProductList(ctx context.Context, limit, offset int) (*[]models.Product, error){
	query := `SELECT id,name,description,price,stock_quantity,created_at FROM products ORDER BY id LIMIT $1 OFFSET $2`

	var products []models.Product

	row,err := repo.DB.Query(ctx,query,limit,offset)

	if err != nil {
		return nil,err
	}

	for row.Next() {
		var product models.Product

		err := row.Scan(&product.ID,&product.Name,&product.Description,&product.Price,&product.StockQuantity,&product.Created_At)

		if err != nil {
			return nil,err
		}

		products = append(products, product)
	}

	defer row.Close()

	return &products,nil

}



func (repo *PostgreDB) ListByCategory(
    ctx context.Context,
    categoryID int64,
    limit, offset int,
) (*[]models.Product, error) {
    query := `
        SELECT id, name, description, price, stock_quantity, created_at
        FROM products
        WHERE category_id = $1
        ORDER BY id
        LIMIT $2 OFFSET $3;
    `

    products := make([]models.Product, 0)

    rows, err := repo.DB.Query(ctx, query, categoryID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var product models.Product

        err := rows.Scan(
            &product.ID,
            &product.Name,
            &product.Description,
            &product.Price,
            &product.StockQuantity,
            &product.Created_At,
        )
        if err != nil {
            return nil, err
        }

        products = append(products, product)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return &products, nil
}
