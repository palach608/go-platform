package repository

import (
	"github.com/google/uuid"
	domainModel "github.com/palach608/go-platform/services/notes/internal/domain/model"
)

type NoteRepository interface {
	Create(note *domainModel.Note) (*domainModel.Note, error)
	Delete(id uuid.UUID) error
	FindAll() ([]*domainModel.Note, error)
	FindByID(id uuid.UUID) (*domainModel.Note, error)
	FindByUserID(userID uuid.UUID) ([]*domainModel.Note, error)
	Update(note *domainModel.Note) (*domainModel.Note, error)
}
