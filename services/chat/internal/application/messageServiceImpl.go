package application

import (
	"context"

	"github.com/google/uuid"
	domainModel "github.com/palach608/go-platform/services/chat/internal/domain/model"
	"github.com/palach608/go-platform/services/chat/internal/domain/service"
	"github.com/palach608/go-platform/services/chat/internal/infrastructure/repository"
)

type MessageServiceImpl struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) service.MessageService {
	return &MessageServiceImpl{repo: repo}
}

func (s *MessageServiceImpl) Create(ctx context.Context, roomID, userID uuid.UUID, body string) (*domainModel.Message, error) {
	message := &domainModel.Message{
		RoomId: roomID,
		UserId: userID,
		Body:   body,
	}

	message, err := s.repo.Create(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageServiceImpl) GetByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*domainModel.Message, error) {
	messages, err := s.repo.FindByRoomID(ctx, roomID, limit, offset)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
