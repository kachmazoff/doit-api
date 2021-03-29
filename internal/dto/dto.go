package dto

import "github.com/kachmazoff/doit-api/internal/model"

type MessageResponse struct {
	Message string `json:"message"`
}

type IdResponse struct {
	Id string `json:"id"`
}

type TokenResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

type Registration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
