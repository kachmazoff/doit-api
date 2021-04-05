package impl

import (
	"github.com/kachmazoff/doit-api/internal/model"
	"github.com/kachmazoff/doit-api/internal/repository"
)

type ChallengesService struct {
	repo             repository.Challenges
	participantsRepo repository.Participants
}

func NewChallengesService(repo repository.Challenges, participantsRepo repository.Participants) *ChallengesService {
	return &ChallengesService{
		repo:             repo,
		participantsRepo: participantsRepo,
	}
}

func (s *ChallengesService) Create(newChallenge model.Challenge) (string, error) {
	return s.repo.Create(newChallenge)
}

func (s *ChallengesService) GetById(id string) (model.Challenge, error) {
	return s.repo.GetById(id)
}

func (s *ChallengesService) GetAll() ([]model.Challenge, error) {
	challenges, err := s.repo.GetAll()
	if err != nil {
		return []model.Challenge{}, err
	}

	for i := 0; i < len(challenges); i++ {
		s.Anonymize(&challenges[i])
	}

	return challenges, nil
}

func (s *ChallengesService) GetAllPublic() ([]model.Challenge, error) {
	challenges, err := s.repo.GetAllPublic()
	if err != nil {
		return []model.Challenge{}, err
	}

	for i := 0; i < len(challenges); i++ {
		s.Anonymize(&challenges[i])
	}

	return challenges, nil
}

func (s *ChallengesService) GetAllOwn(userId string) ([]model.Challenge, error) {
	return s.repo.GetAllOwn(userId)
}

func (s *ChallengesService) Anonymize(challenge *model.Challenge) bool {
	isAnonym := false

	if !challenge.ShowAuthor {
		isAnonym = true
		challenge.AuthorId = ""
	}

	return isAnonym
}

func (s *ChallengesService) EnrichWithUserParticipant(challenge *model.Challenge, userId string) {
	participant, err := s.participantsRepo.GetByUserInChallenge(userId, challenge.Id)
	if err == nil {
		challenge.Participant = &participant
	}
}

func (s *ChallengesService) EnrichAllWithUserParticipant(challenges *[]model.Challenge, userId string) {
	for i := 0; i < len(*challenges); i++ {
		s.EnrichWithUserParticipant(&(*challenges)[i], userId)
	}
}
