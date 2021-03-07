package mysql

import (
	"fmt"

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
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=? LIMIT 1", suggestionsTable)

	var suggestion model.Suggestion
	if err := r.db.Get(&suggestion, query, id); err != nil {
		return model.Suggestion{}, err
	}

	return suggestion, nil
}
