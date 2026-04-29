package service

import (
	"context"

	"github.com/google/uuid"
	domainModel "github.com/palach608/go-platform/services/chat/internal/domain/model"
)

type RoomService interface {
	Create(ctx context.Context, name, description string, ownerID uuid.UUID) (*domainModel.Room, error)
	GetAll(ctx context.Context) ([]*domainModel.Room, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domainModel.Room, error)
	Delete(ctx context.Context, id uuid.UUID, requesterID uuid.UUID) error
}
