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

func (u *ChallengesService) Create(newChallenge model.Challenge) (string, error) {
	const rootUserId string = "e3921905-9eb8-468a-8434-6b3b992a1987"
	newChallenge.AuthorId = rootUserId

	return u.repo.Create(newChallenge)
}

func (u *ChallengesService) GetAll() ([]model.Challenge, error) {
	return u.repo.GetAll()
}

func (s *ChallengesService) Anonymize(challenge *model.Challenge) bool {
	isAnonym := false

	if !challenge.ShowAuthor {
		isAnonym = true
		challenge.AuthorId = ""
	}

	return isAnonym
}
