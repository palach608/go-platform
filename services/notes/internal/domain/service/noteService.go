package service

import (
	apiModel "github.com/SilverName608/go-notes/internal/api/model"
	domainModel "github.com/SilverName608/go-notes/internal/domain/model"
	"github.com/google/uuid"
)

type NoteService interface {
	Create(userID uuid.UUID, req *apiModel.CreateNoteRequest) (*domainModel.Note, error)
	GetAll() ([]*domainModel.Note, error)
	GetByID(id uuid.UUID) (*domainModel.Note, error)
	Update(id uuid.UUID, req *apiModel.UpdateNoteRequest) (*domainModel.Note, error)
	Delete(id uuid.UUID) error
}
