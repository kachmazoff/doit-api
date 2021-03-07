package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/model"
)

type ParticipantsMysqlRepo struct {
	db *sqlx.DB
}

func NewParticipantsMysqlRepo(db *sqlx.DB) *ParticipantsMysqlRepo {
	return &ParticipantsMysqlRepo{db: db}
}

func (r *ParticipantsMysqlRepo) GetById(id string) (model.Participant, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=? LIMIT 1", participantsTable)

	var participant model.Participant
	if err := r.db.Get(&participant, query, id); err != nil {
		return model.Participant{}, err
	}

	return participant, nil
}
