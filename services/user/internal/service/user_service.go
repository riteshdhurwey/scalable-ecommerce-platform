package service

import (
	"context"
	"fmt"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/dto"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/jwt"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/model"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/repository"
)

type UserService struct {
	repo *repository.UserRepo
}

func InitUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {

	if !isValidInput(*req) {
		return nil, ErrInvalidInput
	}

	emailErr := isValidEmail(req.Email)
	if emailErr != nil {
		return nil, emailErr
	}

	passErr := validatePassword(req.Password)

	if passErr != nil {
		return nil, passErr
	}

	userIDExists, err := s.repo.UserNameExists(ctx, req.UserName)

	if err != nil {
		return nil, err
	}
	if userIDExists {
		return nil, ErrUserIDRegistered
	}

	emailExists, err := s.repo.EmailExists(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if emailExists {
		return nil, ErrEmailRegistered
	}

	passwordHash, err := HashPassword(req.Password)

	if err != nil {
		return nil, fmt.Errorf("Password Hash %w", err)
	}

	user := model.User{UserName: req.UserName, Name: req.Name, Email: req.Email, PasswordHash: passwordHash}

	err = s.repo.Create(ctx, &user)

	if err != nil {
		return nil, err
	}

	userResponce := dto.UserResponse{UserName: req.UserName, Name: req.Name, Email: req.Email}

	return &userResponce, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*string, error) {
	if err := isValidEmail(req.Email); err != nil {
		return nil, err
	}

	userExists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, ErrUserNotFound
	}

	passwordHash, userName, err := s.repo.GetPasswordHash(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Println("Password verified:", varifyPassword(req.Password, passwordHash))
	if !varifyPassword(req.Password, passwordHash) {
		return nil, ErrInvalidCredentials
	}

	req.UserName = userName

	authToken, err := jwt.GenerateAccessToken(req)
	if err != nil {
		return nil, err
	}

	return &authToken, nil
}

func (s *UserService) GetProfile(ctx context.Context, userName string) (*dto.UserResponse, error) {

	userExist,err := s.repo.UserNameExists(ctx,userName)

	if err != nil {
		return nil,err
	}

	if !userExist {
		return nil,ErrUserNotFound
	}

	user, err := s.repo.GetByUserName(ctx, userName)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *UserService) DeleteUser(ctx context.Context,userName string)(bool,error){

	userExists,err := s.repo.UserNameExists(ctx,userName)

	if err != nil {
		return false,err
	}

	if !userExists {
		return false,ErrUserNotFound
	}

	userDeleted,err := s.repo.DeleteUser(ctx,userName)

	if err != nil {
		return false,err
	}

	if userDeleted == false {
		return false,fmt.Errorf("failed to delete user %s",userName)
	}

	return true,nil
}