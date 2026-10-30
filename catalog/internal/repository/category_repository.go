package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
)


type PostgreDB struct {
	DB *pgxpool.Pool
}

func NewPostgreRepo(DB *pgxpool.Pool) *PostgreDB{
	return &PostgreDB{DB: DB}
}

func (repo *PostgreDB) CreateCategory(ctx context.Context,c models.Category) (int64,error){
	query := `INSERT INTO categories ( name,description) VALUES ($1,$2) RETURNING id;`

	var categoryID int64

	err := repo.DB.QueryRow(ctx,query,c.Name,c.Description).Scan(&categoryID)

	if err != nil {
		return 0,err
	}

	return categoryID,nil
}



func (repo *PostgreDB) GetCategories(ctx context.Context) ([]models.Category, error) {
    query := `
        SELECT id, name, description
        FROM categories
        ORDER BY id
    `

	rows, err := repo.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var response []models.Category

	for rows.Next() {
		var category models.Category

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		response = append(response, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate error: %w", err)
	}

	return response, nil
}
