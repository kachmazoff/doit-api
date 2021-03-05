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
