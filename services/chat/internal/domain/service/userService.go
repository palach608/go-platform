package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/palach608/go-platform/services/chat/internal/domain/model"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (token string, err error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}
