package service

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/dto"
)

func isValidInput(req dto.RegisterRequest) bool {
	if req.Name == "" || req.Email == "" || req.Password == "" || req.UserName == "" {
		return false
	}
	return true
}

func isValidEmail(email string) error {
	email = strings.TrimSpace(email)

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return fmt.Errorf("invalid email %w", err)
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 || len(password) > 72 {
		return fmt.Errorf("password should contain min 8 charaters")
	}

	return nil
}
