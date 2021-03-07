package impl

import (
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

func (s *NotesService) GetById(id string) (model.Note, error) {
	note, err := s.repo.GetById(id)

	if err != nil {
		return model.Note{}, err
	}

	if *note.Anonymous {
		note.AuthorId = ""
	}

	return note, nil
}
