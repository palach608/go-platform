package application

import (
	"context"

	"github.com/palach608/go-platform/services/auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo domain.UserRepository
}

func NewAuthService(repo domain.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.repo.Create(ctx, user)
}

func (s *AuthService) UpdateUser(ctx context.Context, user *domain.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return s.repo.Update(ctx, user)
}

func (s *AuthService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
