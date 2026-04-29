package repository

import (
	"context"

	domainModel "github.com/SilverName608/go-chat/internal/domain/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domainModel.User) (*domainModel.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domainModel.User, error)
	FindByEmail(ctx context.Context, email string) (*domainModel.User, error)
	Update(ctx context.Context, user *domainModel.User) (*domainModel.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
