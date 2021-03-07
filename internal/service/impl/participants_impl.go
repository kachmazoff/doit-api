package impl

import (
	"errors"

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

func (u *ParticipantsService) GetById(id string) (model.Participant, error) {
	return model.Participant{}, errors.New("Not implemented")
}
