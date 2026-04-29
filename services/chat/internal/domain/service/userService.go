package service

import (
	"context"

	"github.com/SilverName608/go-chat/internal/domain/model"
	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (token string, err error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}
