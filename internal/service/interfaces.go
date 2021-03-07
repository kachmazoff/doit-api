package service

import "github.com/kachmazoff/doit-api/internal/model"

type Users interface {
	Create(user model.User) (string, error)
	ConfirmAccount(userId string) error
	GetByUsername(username string) (model.User, error)
}

type Challenges interface {
	Create(challenge model.Challenge) (string, error)
	GetAll() ([]model.Challenge, error)
}

type Timeline interface {
	GetAll() ([]model.TimelineItem, error)
	GetCommon() ([]model.TimelineItem, error)
}

type Participants interface {
	GetById(id string) (model.Participant, error)
}

type Notes interface {
	GetById(id string) (model.Note, error)
}

type Suggestions interface {
	GetById(id string) (model.Suggestion, error)
}
