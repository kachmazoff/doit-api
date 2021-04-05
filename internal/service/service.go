package service

import (
	"github.com/kachmazoff/doit-api/internal/mailing"
	"github.com/kachmazoff/doit-api/internal/repository"
	"github.com/kachmazoff/doit-api/internal/service/impl"
)

type Services struct {
	Users
	Challenges
	Participants
	Notes
	Suggestions
	Timeline
	Followers
}

func NewServices(r *repository.Repositories, sender mailing.Sender) *Services {
	users := impl.NewUsersService(r.Users, sender)
	challenges := impl.NewChallengesService(r.Challenges, r.Participants)
	participants := impl.NewParticipantsService(r.Participants)
	notes := impl.NewNotesService(r.Notes)
	suggestions := impl.NewSuggestionsService(r.Suggestions)

	timeline := impl.NewTimelineService(r.Timeline, *challenges, *notes, *participants, *suggestions)

	followers := impl.NewFollowersService(r.Followers)

	return &Services{
		Users:        users,
		Challenges:   challenges,
		Participants: participants,
		Notes:        notes,
		Suggestions:  suggestions,
		Timeline:     timeline,
		Followers:    followers,
	}
}
