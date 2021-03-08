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

func (r *SuggestionsMysqlRepo) GetForParticipant(participantId string) ([]model.Suggestion, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE participant_id=?", suggestionsTable)

	return r.selectSuggestions(query, participantId)
}

func (r *SuggestionsMysqlRepo) GetByAuthor(authorId string, onlyPublic bool) ([]model.Suggestion, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE author_id=?", suggestionsTable)
	if onlyPublic {
		query += " AND anonymous=true"
	}

	return r.selectSuggestions(query, authorId)
}

func (r *SuggestionsMysqlRepo) selectSuggestions(query string, args ...interface{}) ([]model.Suggestion, error) {
	var suggestions []model.Suggestion
	if err := r.db.Select(&suggestions, query, args); err != nil {
		return []model.Suggestion{}, err
	}

	return suggestions, nil
}
