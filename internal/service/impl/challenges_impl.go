package impl

import (
	"github.com/kachmazoff/doit-api/internal/model"
	"github.com/kachmazoff/doit-api/internal/repository"
)

type ChallengesService struct {
	repo repository.Challenges
}

func NewChallengesService(repo repository.Challenges) *ChallengesService {
	return &ChallengesService{
		repo: repo,
	}
}

func (s *ChallengesService) Create(newChallenge model.Challenge) (string, error) {
	return s.repo.Create(newChallenge)
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

func (s *ChallengesService) Anonymize(challenge *model.Challenge) bool {
	isAnonym := false

	if !challenge.ShowAuthor {
		isAnonym = true
		challenge.AuthorId = ""
	}

	return isAnonym
}
