package mysql

import (
	"errors"

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
	return model.Note{}, errors.New("Not implemented")
}
