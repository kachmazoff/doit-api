package model

import "time"

type TimelineItem struct {
	Index   uint      `json:"index" db:"index"`
	Created time.Time `json:"created" db:"created"`
	Type    string    `json:"type" db:"type"`

	UserId      string `json:"-" db:"user_id"`
	ChallengeId string `json:"-" db:"challenge_id"`

	ParticipantId *string `json:"-" db:"participant_id"`
	NoteId        *string `json:"-" db:"note_id"`
	SuggestionId  *string `json:"-" db:"suggestion_id"`

	*User        `json:"user,omitempty" db:"user"`
	Challenge    `json:"challenge" db:"challenge"`
	*Participant `json:"participant,omitempty"`
	*Note        `json:"note,omitempty"`
	*Suggestion  `json:"suggestion,omitempty"`
}

type TimelineFilters struct {
	RequestAuthor string
	UserId        string
	TimelineType  string // common, own, subs
	EventTypes    []string
	ParticipantId string
	ChallengeId   string
	Limit         int
	LastIndex     int
	OrderType     string
}
