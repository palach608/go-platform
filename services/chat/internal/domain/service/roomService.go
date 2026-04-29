package service

import (
	"context"

	domainModel "github.com/SilverName608/go-chat/internal/domain/model"
	"github.com/google/uuid"
)

type RoomService interface {
	Create(ctx context.Context, name, description string, ownerID uuid.UUID) (*domainModel.Room, error)
	GetAll(ctx context.Context) ([]*domainModel.Room, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domainModel.Room, error)
	Delete(ctx context.Context, id uuid.UUID, requesterID uuid.UUID) error
}
