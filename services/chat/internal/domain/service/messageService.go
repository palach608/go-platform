package service

import (
	"context"

	"github.com/google/uuid"

	domainModel "github.com/SilverName608/go-chat/internal/domain/model"
)

type MessageService interface {
	Create(ctx context.Context, roomID, userID uuid.UUID, body string) (*domainModel.Message, error)
	GetByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*domainModel.Message, error)
}
