package application

import (
	"context"

	"github.com/google/uuid"
	kafkapkg "github.com/palach608/go-platform/pkg/broker"
	apiModel "github.com/palach608/go-platform/services/notes/internal/api/model"
	domainModel "github.com/palach608/go-platform/services/notes/internal/domain/model"
	"github.com/palach608/go-platform/services/notes/internal/domain/service"
	"github.com/palach608/go-platform/services/notes/internal/infrastructure/repository"
)

type NoteServiceImpl struct {
	repo     repository.NoteRepository
	producer *kafkapkg.Producer
}

func NewNoteService(repo repository.NoteRepository, producer *kafkapkg.Producer) service.NoteService {
	return &NoteServiceImpl{repo: repo, producer: producer}
}

func (s *NoteServiceImpl) Create(ctx context.Context, userID uuid.UUID, username string, req *apiModel.CreateNoteRequest) (*domainModel.Note, error) {
	note := &domainModel.Note{
		UserID: userID,
		Title:  req.Title,
		Body:   req.Body,
	}

	note, err := s.repo.Create(note)
	if err != nil {
		return nil, err
	}

	event := kafkapkg.NoteCreatedEvent{
		Type:     "note.created",
		UserID:   userID.String(),
		Username: username,
		NoteID:   note.ID.String(),
		Title:    note.Title,
	}

	_ = s.producer.Publish(ctx, note.ID.String(), event)

	return note, nil
}

func (s *NoteServiceImpl) GetAll() ([]*domainModel.Note, error) {
	return s.repo.FindAll()
}

func (s *NoteServiceImpl) GetByID(id uuid.UUID) (*domainModel.Note, error) {
	return s.repo.FindByID(id)
}

func (s *NoteServiceImpl) Update(id uuid.UUID, req *apiModel.UpdateNoteRequest) (*domainModel.Note, error) {
	note, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	note.Title = req.Title
	note.Body = req.Body
	return s.repo.Update(note)
}

func (s *NoteServiceImpl) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
