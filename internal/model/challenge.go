package model

import "time"

type Challenge struct {
	Id               string    `json:"id" db:"id"`
	Created          time.Time `json:"created" db:"created"`
	AuthorId         string    `json:"author_id,omitempty" db:"author_id"`
	ShowAuthor       bool      `json:"show_author" db:"show_author"`
	Title            string    `json:"title" db:"title"`
	Body             *string   `json:"body" db:"body"`
	VisibleType      string    `json:"visible_type" db:"visible_type"`
	ParticipantsType string    `json:"participants_type" db:"participants_type"`
}
