package model

import "time"

type Note struct {
	Id            string    `json:"id" db:"id"`
	Created       time.Time `json:"created" db:"created"`
	ChallengeId   string    `json:"challenge_id" db:"challenge_id"`
	Type          string    `json:"type" db:"type"`
	Body          string    `json:"body" db:"body"`
	ParticipantId string    `json:"participant_id" db:"participant_id"`
	AuthorId      string    `json:"author_id" db:"author_id"`
}
