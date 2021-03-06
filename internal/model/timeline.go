package model

import "time"

type TimelineItem struct {
	Index   uint      `json:"index" db:"index"`
	Created time.Time `json:"created" db:"created"`
	UserId  string    `json:"user_id" db:"user_id"`
	Type    string    `json:"type" db:"type"`

	*User        `json:"user" db:"user"`
	Challenge    `json:"challenge" db:"challenge"`
	*Participant `json:"participant" db:"participant"`
	*Note        `json:"note" db:"note"`
	*Suggestion  `json:"suggestion" db:"suggestion"`
}
