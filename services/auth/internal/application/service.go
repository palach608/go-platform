package application

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/palach608/go-platform/services/auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo      domain.UserRepository
	jwtSecret string
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

func NewAuthService(repo domain.UserRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	return &AuthService{
		repo:      repo,
		jwtSecret: secret,
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	existing, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return ErrUserAlreadyExists
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &domain.User{Username: username, Email: email, Password: string(hash)}

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
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
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

func (s *AuthService) DeleteUser(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
