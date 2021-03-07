package impl

import (
	"errors"

	"github.com/kachmazoff/doit-api/internal/model"
	"github.com/kachmazoff/doit-api/internal/repository"
)

type NotesService struct {
	repo repository.Notes
}

func NewNotesService(repo repository.Notes) *NotesService {
	return &NotesService{
		repo: repo,
	}
}

func (u *NotesService) GetById(id string) (model.Note, error) {
	return model.Note{}, errors.New("Not implemented")
}
