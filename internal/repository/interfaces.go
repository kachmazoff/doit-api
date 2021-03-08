package repository

import "github.com/kachmazoff/doit-api/internal/model"

type Users interface {
	Create(model.User) (string, error)
	// GetAll() ([]model.User, error)
	GetByUsername(username string) (model.User, error)
	SetStatus(id string, status string) error
	GetEmailById(id string) (string, error)
}

type Challenges interface {
	Create(model.Challenge) (string, error)
	// TODO: Remove later. Temp method
	GetAll() ([]model.Challenge, error)
}

type Timeline interface {
	CreateChallenge() error
	// TODO: Remove later. Temp method
	GetAll() ([]model.TimelineItem, error)
	GetCommon() ([]model.TimelineItem, error)
	GetForUser(userId string) ([]model.TimelineItem, error)
	GetUserOwn(userId string) ([]model.TimelineItem, error)
}

type Participants interface {
	Create(participant model.Participant) (string, error)
	GetById(id string) (model.Participant, error)
	GetParticipationsOfUser(userId string, onlyPublic, onlyActive bool) ([]model.Participant, error)
	GetParticipantsInChallenge(challengeId string, onlyPublic, onlyActive bool) ([]model.Participant, error)
}

type Notes interface {
	Create(note model.Note) (string, error)
	GetById(id string) (model.Note, error)
	GetNotesOfParticipant(participantId string) ([]model.Note, error)
}

type Suggestions interface {
	Create(suggestion model.Suggestion) (string, error)
	GetById(id string) (model.Suggestion, error)
	GetForParticipant(participantId string) ([]model.Suggestion, error)
	GetByAuthor(authorId string, onlyPublic bool) ([]model.Suggestion, error)
}

type Followers interface {
	Subscribe(fromId string, toId string) error
	Unsubscribe(fromId string, toId string) error
	GetFollowersIds(userId string) ([]string, error)
	GetFollowedIds(userId string) ([]string, error)
	GetFollowers(userId string) ([]model.User, error)
	GetFollowees(userId string) ([]model.User, error)
}
