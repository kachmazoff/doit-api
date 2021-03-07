package model

import "time"

type Note struct {
	Id            string    `json:"id" db:"id"`
	Created       time.Time `json:"created" db:"created"`
	Type          string    `json:"type" db:"type"`
	Body          string    `json:"body" db:"body"`
	ParticipantId string    `json:"participant_id" db:"participant_id"`
	AuthorId      string    `json:"author_id" db:"author_id"`
	Anonymous     *bool     `json:"-" db:"anonymous"`
}
