package model

import "time"

type Suggestion struct {
	Id            string    `json:"id" db:"id"`
	Created       time.Time `json:"created" db:"created"`
	ChallengeId   string    `json:"challenge_id" db:"challenge_id"`
	Body          string    `json:"body" db:"body"`
	ParticipantId string    `json:"participant_id" db:"participant_id"`
	AuthorId      string    `json:"author_id" db:"author_id"`
	Anonymous     bool      `json:"anonymous" db:"anonymous"`
	CreatedNoteId *string   `json:"created_note_id" db:"created_note_id"`
}
