package service

import "github.com/kachmazoff/doit-api/internal/model"

type Users interface {
	Create(user model.User) (string, error)
	ConfirmAccount(userId string) error
	GetAll() ([]model.User, error)
	GetByUsername(username string) (model.User, error)
	GetByEmail(email string) (model.User, error)
	GetIdByUsername(username string) (string, error)
}

type Challenges interface {
	Create(challenge model.Challenge) (string, error)
	GetById(id string) (model.Challenge, error)
	GetAll() ([]model.Challenge, error)
	GetAllPublic() ([]model.Challenge, error)
	GetAllOwn(userId string) ([]model.Challenge, error)
	Anonymize(*model.Challenge) bool
	EnrichWithUserParticipant(challenge *model.Challenge, userId string)
	EnrichAllWithUserParticipant(challenges *[]model.Challenge, userId string)
}

type Timeline interface {
	GetWithFilters(filters model.TimelineFilters) ([]model.TimelineItem, error)
	GetAll() ([]model.TimelineItem, error)
	GetCommon() ([]model.TimelineItem, error)
	GetForUser(userId string) ([]model.TimelineItem, error)
	GetUserOwn(userId string) ([]model.TimelineItem, error)
	Anonymize(*[]model.TimelineItem) bool
	AnonymizeItem(*model.TimelineItem) bool
}

type Participants interface {
	Create(participant model.Participant) (string, error)
	GetByIdUNSAFE(id string) (model.Participant, error)
	GetById(id string) (model.Participant, error)
	GetParticipationsOfUser(userId string, onlyPublic, onlyActive bool) ([]model.Participant, error)
	GetParticipantsInChallenge(challengeId string, onlyPublic, onlyActive bool) ([]model.Participant, error)
	HasRootAccess(participantId, userId string) bool
	IsPublic(participantId string) bool
	Anonymize(participant *model.Participant, userId string) bool
}

type Notes interface {
	Create(note model.Note) (string, error)
	GetById(id string) (model.Note, error)
	GetNotesOfParticipant(participantId string, needAnonymize bool) ([]model.Note, error)
	Anonymize(*model.Note) bool
}

type Suggestions interface {
	Create(suggestion model.Suggestion) (string, error)
	GetById(id string) (model.Suggestion, error)
	GetForParticipant(participantId string) ([]model.Suggestion, error)
	// TODO: подумать над тем, как корректно объявить api-эндпоинты для этого и следующего метода
	GetByAuthor(authorId string, onlyPublic bool) ([]model.Suggestion, error)
	GetForUser(userId string) ([]model.Suggestion, error)
	Anonymize(*model.Suggestion) bool
}

type Followers interface {
	Subscribe(fromId, toId string) error
	Unsubscribe(fromId, toId string) error

	GetFollowersIds(userId string) ([]string, error)
	GetFollowedIds(userId string) ([]string, error)

	GetFollowers(userId string) ([]model.User, error)
	GetFollowees(userId string) ([]model.User, error)

	ExistsFromTo(fromId, toId string) (bool, error)
}
