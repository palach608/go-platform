package repository

import (
	"context"

	domainModel "github.com/SilverName608/go-chat/internal/domain/model"

	"github.com/google/uuid"
)

type MessageRepository interface {
	Create(ctx context.Context, message *domainModel.Message) (*domainModel.Message, error)
	FindByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*domainModel.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
