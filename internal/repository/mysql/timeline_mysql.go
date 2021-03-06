package mysql

import "github.com/jmoiron/sqlx"

type TimelineMysqlRepo struct {
	db *sqlx.DB
}

func NewTimelineMysqlRepo(db *sqlx.DB) *TimelineMysqlRepo {
	return &TimelineMysqlRepo{db: db}
}
