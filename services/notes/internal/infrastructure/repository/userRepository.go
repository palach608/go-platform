package repository

import (
	domainModel "github.com/SilverName608/go-notes/internal/domain/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domainModel.User) (*domainModel.User, error)
	Delete(id uuid.UUID) error
	FindByEmail(email string) (*domainModel.User, error)
	FindByID(id uuid.UUID) (*domainModel.User, error)
}
