package application

import (
	"context"
	"errors"
	"time"

	domainModel "github.com/SilverName608/go-chat/internal/domain/model"
	"github.com/SilverName608/go-chat/internal/domain/service"
	"github.com/SilverName608/go-chat/internal/infrastructure/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) service.UserService {
	return &UserServiceImpl{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserServiceImpl) Register(ctx context.Context, username, email, password string) (*domainModel.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domainModel.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	user, err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := t.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*domainModel.User, error) {
	return s.repo.FindByID(ctx, id)
}
