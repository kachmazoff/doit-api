package impl

import (
	"github.com/kachmazoff/doit-api/internal/model"
	"github.com/kachmazoff/doit-api/internal/repository"
)

type FollowersService struct {
	repo repository.Followers
}

func NewFollowersService(repo repository.Followers) *FollowersService {
	return &FollowersService{
		repo: repo,
	}
}

func (s *FollowersService) Subscribe(fromId, toId string) error {
	return s.repo.Subscribe(fromId, toId)
}

func (s *FollowersService) Unsubscribe(fromId, toId string) error {
	return s.repo.Unsubscribe(fromId, toId)
}

func (s *FollowersService) GetFollowersIds(userId string) ([]string, error) {
	return s.repo.GetFollowersIds(userId)
}

func (s *FollowersService) GetFollowedIds(userId string) ([]string, error) {
	return s.repo.GetFollowedIds(userId)
}

func (s *FollowersService) GetFollowers(userId string) ([]model.User, error) {
	return s.repo.GetFollowers(userId)
}

func (s *FollowersService) GetFollowees(userId string) ([]model.User, error) {
	return s.repo.GetFollowees(userId)
}
