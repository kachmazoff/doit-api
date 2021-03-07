package impl

import (
	"github.com/kachmazoff/doit-api/internal/model"
	"github.com/kachmazoff/doit-api/internal/repository"
)

type ParticipantsService struct {
	repo repository.Participants
}

func NewParticipantsService(repo repository.Participants) *ParticipantsService {
	return &ParticipantsService{
		repo: repo,
	}
}

func (s *ParticipantsService) GetById(id string) (model.Participant, error) {
	participant, err := s.repo.GetById(id)

	if err != nil {
		return model.Participant{}, err
	}

	s.Anonymize(&participant)

	return participant, nil
}

func (s *ParticipantsService) Anonymize(participant *model.Participant) bool {
	isAnonym := false

	if participant.Anonymous {
		isAnonym = true
		participant.UserId = ""
		participant.TeamId = nil
	}

	return isAnonym
}
