package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/dto"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/model"
)

type UserRepo struct {
	DB *pgxpool.Pool
}

func InitUserRepo(DB *pgxpool.Pool) *UserRepo {
	return &UserRepo{DB: DB}
}

func (repo *UserRepo) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (user_name,name,email,password_hash) VALUES ($1,$2,$3,$4);`

	err := repo.DB.QueryRow(ctx, query, user.UserName, user.Name, user.Email, user.PasswordHash)

	if err != nil {
		return ctx.Err()
	}

	return nil

}

func (repo *UserRepo) GetByUserName(ctx context.Context, userName string) (*dto.UserResponse, error) {
	query := `SELECT user_name,name,email,created_at FROM users WHERE user_name = $1;`

	var user dto.UserResponse

	err := repo.DB.QueryRow(ctx, query, userName).Scan(&user.UserName, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepo) GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	query := `SELECT user_name,name,email,created_at FROM users WHERE email = $1;`

	var user dto.UserResponse

	err := repo.DB.QueryRow(ctx, query, email).Scan(&user.UserName, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepo) Update(ctx context.Context, user *model.User) (*dto.UserResponse, error) {
	query := `UPDATE users SET name = $1, email = $2 WHERE user_name = $3 RETURNING created_at`

	var Created_At time.Time
	err := repo.DB.QueryRow(ctx, query, user.Name, user.Email, user.UserName).Scan(&Created_At)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	response := dto.UserResponse{UserName: user.UserName, Name: user.Name, Email: user.Email, CreatedAt: Created_At}

	return &response, nil

}

func (repo *UserRepo) UserNameExists(ctx context.Context, userName string) (bool, error) {
	query := `SELECT EXISTS (
			  SELECT 1 
			  FROM users 
			  WHERE user_name = $1
	);`

	var exists bool

	err := repo.DB.QueryRow(ctx, query, userName).Scan(&exists)

	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	return true, nil
}

func (repo *UserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);`

	var exists bool

	err := repo.DB.QueryRow(ctx, query, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	return true, nil
}

func (repo *UserRepo) GetPasswordHash(ctx context.Context, req dto.LoginRequest) (string, string, error) {
	var query string
	var args []any

	switch {
	case req.UserName != "" && req.Email != "":
		query = `SELECT password_hash, user_name
		         FROM users
		         WHERE user_name = $1 AND email = $2
		         LIMIT 1`
		args = []any{req.UserName, req.Email}

	case req.UserName != "":
		query = `SELECT password_hash, user_name
		         FROM users
		         WHERE user_name = $1
		         LIMIT 1`
		args = []any{req.UserName}

	case req.Email != "":
		query = `SELECT password_hash, user_name
		         FROM users
		         WHERE email = $1
		         LIMIT 1`
		args = []any{req.Email}

	default:
		return "", "", errors.New("user ID or email is required")
	}

	var passwordHash string
	var userID string

	err := repo.DB.QueryRow(ctx, query, args...).Scan(
		&passwordHash,
		&userID,
	)
	if err != nil {
		return "", "", err
	}

	return passwordHash, userID, nil
}


func (repo *UserRepo) DeleteUser(ctx context.Context,userName string)(bool,error){
	query := `DELETE FROM users WHERE user_name = $1;`
	result,err := repo.DB.Exec(ctx,query,userName)

	if err != nil {
		return false,err
	}

	if result.RowsAffected() == 0 {
		return false,nil
	}

	return true,nil
}