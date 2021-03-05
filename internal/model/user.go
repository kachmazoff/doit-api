package model

import "time"

type User struct {
	Id       string    `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	Created  time.Time `json:"created" db:"created"`
	Password string    `json:"-" db:"password"`
	Status   string    `json:"status" db:"account_status"`
}
