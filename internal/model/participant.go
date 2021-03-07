package model

import "time"

type Participant struct {
	Id          string     `json:"id" db:"id"`
	Created     time.Time  `json:"created" db:"created"`
	ChallengeId string     `json:"challenge_id" db:"challenge_id"`
	UserId      string     `json:"user_id,omitempty" db:"user_id"`
	TeamId      *string    `json:"team_id,omitempty" db:"team_id"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	Status      string     `json:"status" db:"status"`
	Anonymous   bool       `json:"anonymous" db:"anonymous"`
	VisibleType string     `json:"visible_type" db:"visible_type"`
}
