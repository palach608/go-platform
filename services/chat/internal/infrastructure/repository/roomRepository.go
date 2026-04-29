package repository

import (
	"context"

	domainModel "github.com/SilverName608/go-chat/internal/domain/model"

	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(ctx context.Context, room *domainModel.Room) (*domainModel.Room, error)
	FindAll(ctx context.Context) ([]*domainModel.Room, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domainModel.Room, error)
	Update(ctx context.Context, room *domainModel.Room) (*domainModel.Room, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
