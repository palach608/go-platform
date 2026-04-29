package service

import (
	"github.com/google/uuid"
	apiModel "github.com/palach608/go-platform/services/notes/internal/api/model"
	domainModel "github.com/palach608/go-platform/services/notes/internal/domain/model"
)

type NoteService interface {
	Create(userID uuid.UUID, req *apiModel.CreateNoteRequest) (*domainModel.Note, error)
	GetAll() ([]*domainModel.Note, error)
	GetByID(id uuid.UUID) (*domainModel.Note, error)
	Update(id uuid.UUID, req *apiModel.UpdateNoteRequest) (*domainModel.Note, error)
	Delete(id uuid.UUID) error
}
