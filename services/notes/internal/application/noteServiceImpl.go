package application

import (
	apiModel "github.com/SilverName608/go-notes/internal/api/model"
	domainModel "github.com/SilverName608/go-notes/internal/domain/model"
	"github.com/SilverName608/go-notes/internal/domain/service"
	"github.com/SilverName608/go-notes/internal/infrastructure/repository"
	"github.com/google/uuid"
)

type NoteServiceImpl struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) service.NoteService {
	return &NoteServiceImpl{repo: repo}
}

func (s *NoteServiceImpl) Create(userID uuid.UUID, req *apiModel.CreateNoteRequest) (*domainModel.Note, error) {
	note := &domainModel.Note{
		UserID: userID,
		Title:  req.Title,
		Body:   req.Body,
	}

	note, err := s.repo.Create(note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteServiceImpl) GetAll() ([]*domainModel.Note, error) {
	notes, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *NoteServiceImpl) GetByID(id uuid.UUID) (*domainModel.Note, error) {
	note, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteServiceImpl) Update(id uuid.UUID, req *apiModel.UpdateNoteRequest) (*domainModel.Note, error) {
	note, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	note.Title = req.Title
	note.Body = req.Body

	note, err = s.repo.Update(note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteServiceImpl) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
