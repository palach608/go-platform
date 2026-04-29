package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	domainModel "github.com/palach608/go-platform/services/chat/internal/domain/model"
	"github.com/palach608/go-platform/services/chat/internal/domain/service"
	"github.com/palach608/go-platform/services/chat/internal/infrastructure/repository"
)

type RoomServiceImpl struct {
	repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) service.RoomService {
	return &RoomServiceImpl{repo: repo}
}

func (s *RoomServiceImpl) Create(ctx context.Context, name, description string, ownerID uuid.UUID) (*domainModel.Room, error) {
	room := &domainModel.Room{
		Name:        name,
		Description: description,
		OwnerId:     ownerID,
	}

	room, err := s.repo.Create(ctx, room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomServiceImpl) GetAll(ctx context.Context) ([]*domainModel.Room, error) {
	rooms, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *RoomServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*domainModel.Room, error) {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomServiceImpl) Delete(ctx context.Context, id uuid.UUID, requesterID uuid.UUID) error {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if room.OwnerId != requesterID {
		return errors.New("forbidden: only owner can delete room")
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
