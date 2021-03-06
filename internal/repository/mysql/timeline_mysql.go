package mysql

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/model"
)

type TimelineMysqlRepo struct {
	db *sqlx.DB
}

func NewTimelineMysqlRepo(db *sqlx.DB) *TimelineMysqlRepo {
	return &TimelineMysqlRepo{db: db}
}

func (r *TimelineMysqlRepo) CreateChallenge() error {
	// query := fmt.Sprintf("INSERT INTO %s () VALUES ()", timelineTable)
	return errors.New("Not implemented")
}

func (r *TimelineMysqlRepo) GetAll() ([]model.TimelineItem, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY ASC", timelineTable)

	var timeline []model.TimelineItem
	if err := r.db.Select(&timeline, query); err != nil {
		println(err.Error())
		return []model.TimelineItem{}, err
	}

	return timeline, nil
}
