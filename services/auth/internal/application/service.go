package application

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/palach608/go-platform/services/auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      domain.UserRepository
	jwtSecret string
}

func NewAuthService(repo domain.UserRepository) *AuthService {
	// Секрет потом прокинем через DI из конфига
	return &AuthService{
		repo:      repo,
		jwtSecret: "your-very-secret-key",
	}
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

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) UpdateUser(ctx context.Context, user *domain.User) error {
	if user.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}
	return s.repo.Update(ctx, user)
}

func (s *AuthService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
