package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/repository/mysql"
)

type Repositories struct {
	Users
	Challenges
	Timeline
	Participants
	Notes
	Suggestions
	Followers
}

func NewMysqlRepos(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users:        mysql.NewUsersMysqlRepo(db),
		Challenges:   mysql.NewChallengesMysqlRepo(db),
		Timeline:     mysql.NewTimelineMysqlRepo(db),
		Participants: mysql.NewParticipantsMysqlRepo(db),
		Notes:        mysql.NewNotesMysqlRepo(db),
		Suggestions:  mysql.NewSuggestionsMysqlRepo(db),
		Followers:    mysql.NewFollowersMysqlRepo(db),
	}
}
