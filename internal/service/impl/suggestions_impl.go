package impl

import (
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

func (s *SuggestionsService) GetById(id string) (model.Suggestion, error) {
	suggestion, err := s.repo.GetById(id)

	if err != nil {
		return model.Suggestion{}, err
	}

	if suggestion.Anonymous {
		suggestion.AuthorId = ""
	}

	return suggestion, nil
}
