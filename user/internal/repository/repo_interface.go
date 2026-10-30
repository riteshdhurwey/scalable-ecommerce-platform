package repository

import (
	"context"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/dto"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByUserID(ctx context.Context, userName string) (*dto.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	Update(ctx context.Context, user *model.User) (*dto.UserResponse, error)
	UserIDExists(ctx context.Context, userName string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	DeleteUser(ctx context.Context,userName string) (bool,error)
}
