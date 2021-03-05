package service

import (
	"github.com/kachmazoff/doit-api/internal/mailing"
	"github.com/kachmazoff/doit-api/internal/repository"
	"github.com/kachmazoff/doit-api/internal/service/impl"
)

type Services struct {
	Users
	Challenges
}

func NewServices(r *repository.Repositories, sender mailing.Sender) *Services {
	return &Services{
		Users:      impl.NewUsersService(r.Users, sender),
		Challenges: impl.NewChallengesService(r.Challenges),
	}
}
