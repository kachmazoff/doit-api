package service

import (
	"github.com/kachmazoff/doit-api/internal/mailing"
	"github.com/kachmazoff/doit-api/internal/repository"
	"github.com/kachmazoff/doit-api/internal/service/impl"
)

type Services struct {
	Users
	Challenges
	Timeline
}

func NewServices(r *repository.Repositories, sender mailing.Sender) *Services {
	users := impl.NewUsersService(r.Users, sender)
	challenges := impl.NewChallengesService(r.Challenges)
	notes := impl.NewNotesService(r.Notes)
	participants := impl.NewParticipantsService(r.Participants)
	suggestions := impl.NewSuggestionsService(r.Suggestions)

	timeline := impl.NewTimelineService(r.Timeline, *challenges, *notes, *participants, *suggestions)

	return &Services{
		Users:      users,
		Challenges: challenges,
		Timeline:   timeline,
	}
}
