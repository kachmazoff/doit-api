package mysql

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/model"
)

type SuggestionsMysqlRepo struct {
	db *sqlx.DB
}

func NewSuggestionsMysqlRepo(db *sqlx.DB) *SuggestionsMysqlRepo {
	return &SuggestionsMysqlRepo{db: db}
}

func (r *SuggestionsMysqlRepo) GetById(id string) (model.Suggestion, error) {
	return model.Suggestion{}, errors.New("Not implemented")
}
