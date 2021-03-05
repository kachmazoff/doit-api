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
