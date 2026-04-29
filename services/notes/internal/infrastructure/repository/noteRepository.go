package repository

import (
	domainModel "github.com/SilverName608/go-notes/internal/domain/model"
	"github.com/google/uuid"
)

type NoteRepository interface {
	Create(note *domainModel.Note) (*domainModel.Note, error)
	Delete(id uuid.UUID) error
	FindAll() ([]*domainModel.Note, error)
	FindByID(id uuid.UUID) (*domainModel.Note, error)
	FindByUserID(userID uuid.UUID) ([]*domainModel.Note, error)
	Update(note *domainModel.Note) (*domainModel.Note, error)
}
