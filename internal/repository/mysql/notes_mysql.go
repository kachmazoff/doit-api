package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/model"
)

type NotesMysqlRepo struct {
	db *sqlx.DB
}

func NewNotesMysqlRepo(db *sqlx.DB) *NotesMysqlRepo {
	return &NotesMysqlRepo{db: db}
}

func (r *NotesMysqlRepo) GetById(id string) (model.Note, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=? LIMIT 1", notesTable)

	var note model.Note
	if err := r.db.Get(&note, query, id); err != nil {
		return model.Note{}, err
	}

	return note, nil
}
