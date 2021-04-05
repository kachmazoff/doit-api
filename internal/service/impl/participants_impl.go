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

func (s *ParticipantsService) Create(participant model.Participant) (string, error) {
	return s.repo.Create(participant)
}

func (s *ParticipantsService) GetByIdUNSAFE(id string) (model.Participant, error) {
	participant, err := s.repo.GetById(id)

	if err != nil {
		return model.Participant{}, err
	}

	return participant, nil
}

func (s *ParticipantsService) GetById(id string) (model.Participant, error) {
	participant, err := s.GetByIdUNSAFE(id)

	if err != nil {
		s.Anonymize(&participant, "")
	}

	return participant, err
}

func (s *ParticipantsService) GetParticipationsOfUser(userId string, onlyPublic, onlyActive bool) ([]model.Participant, error) {
	return s.repo.GetParticipationsOfUser(userId, onlyPublic, onlyActive)
}

func (s *ParticipantsService) GetParticipantsInChallenge(challengeId string, onlyPublic, onlyActive bool) ([]model.Participant, error) {
	participants, err := s.repo.GetParticipantsInChallenge(challengeId, onlyPublic, onlyActive)

	if err != nil {
		return []model.Participant{}, err
	}

	for i := 0; i < len(participants); i++ {
		s.Anonymize(&participants[i], "")
	}

	return participants, nil
}

func (s *ParticipantsService) HasRootAccess(participantId, userId string) bool {
	participant, err := s.repo.GetById(participantId)
	if err != nil {
		// TODO: handle error
		return false
	}
	// TODO: check userId in team
	return participant.UserId == userId
}

func (s *ParticipantsService) IsPublic(participantId string) bool {
	participant, err := s.GetById(participantId)
	return err == nil && participant.VisibleType == "public"
}

func (s *ParticipantsService) Anonymize(participant *model.Participant, userId string) bool {
	isAnonym := false

	if participant.Anonymous && participant.UserId != userId {
		isAnonym = true
		participant.UserId = ""
		participant.TeamId = nil
	}

	if participant.Challenge != nil && participant.Challenge.ShowAuthor == false && participant.Challenge.AuthorId != userId {
		participant.Challenge.AuthorId = ""
	}

	return isAnonym
}
