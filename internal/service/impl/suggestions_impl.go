package impl

import (
	"errors"

	"github.com/kachmazoff/doit-api/internal/model"
	"github.com/kachmazoff/doit-api/internal/repository"
)

type SuggestionsService struct {
	repo repository.Suggestions
}

func NewSuggestionsService(repo repository.Suggestions) *SuggestionsService {
	return &SuggestionsService{
		repo: repo,
	}
}

func (u *SuggestionsService) GetById(id string) (model.Suggestion, error) {
	return model.Suggestion{}, errors.New("Not implemented")
}
