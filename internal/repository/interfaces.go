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
}

type Participants interface {
	Create(participant model.Participant) (string, error)
	GetById(id string) (model.Participant, error)
}

type Notes interface {
	Create(note model.Note) (string, error)
	GetById(id string) (model.Note, error)
}

type Suggestions interface {
	GetById(id string) (model.Suggestion, error)
}
