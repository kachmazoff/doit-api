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

func (s *SuggestionsService) Create(suggestion model.Suggestion) (string, error) {
	return s.repo.Create(suggestion)
}

func (s *SuggestionsService) GetById(id string) (model.Suggestion, error) {
	suggestion, err := s.repo.GetById(id)

	if err != nil {
		return model.Suggestion{}, err
	}

	s.Anonymize(&suggestion)

	return suggestion, nil
}

func (s *SuggestionsService) GetForParticipant(participantId string) ([]model.Suggestion, error) {
	return s.defaultAnonymizedResponse(s.repo.GetForParticipant(participantId))
}

func (s *SuggestionsService) GetByAuthor(authorId string, onlyPublic bool) ([]model.Suggestion, error) {
	return s.repo.GetByAuthor(authorId, onlyPublic)
}

func (s *SuggestionsService) GetForUser(userId string) ([]model.Suggestion, error) {
	return s.defaultAnonymizedResponse(s.repo.GetForUser(userId))
}

func (s *SuggestionsService) Anonymize(suggestion *model.Suggestion) bool {
	isAnonym := false

	if suggestion.Anonymous {
		isAnonym = true
		suggestion.AuthorId = ""
	}

	return isAnonym
}

func (s *SuggestionsService) defaultAnonymizedResponse(suggestions []model.Suggestion, err error) ([]model.Suggestion, error) {
	if err != nil {
		return []model.Suggestion{}, err
	}

	for i := 0; i < len(suggestions); i++ {
		s.Anonymize(&suggestions[i])
	}

	return suggestions, nil
}
