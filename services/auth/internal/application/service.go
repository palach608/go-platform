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
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.repo.Create(ctx, user)
}
