package mysql

import (
	"fmt"

	"github.com/google/uuid"
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
	// query := fmt.Sprintf("SELECT * FROM %s WHERE id=? LIMIT 1", notesTable, participanparticipantsTable)
	query := fmt.Sprintf("SELECT n.*, p.anonymous AS `anonymous` FROM (SELECT * FROM %s WHERE id=? LIMIT 1) AS n LEFT JOIN %s AS p ON n.participant_id=p.id", notesTable, participantsTable)

	var note model.Note
	if err := r.db.Get(&note, query, id); err != nil {
		return model.Note{}, err
	}

	return note, nil
}

func (r *NotesMysqlRepo) Create(note model.Note) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s (id, type, body, participant_id, author_id) values (?, ?, ?, ?, ?)", notesTable)
	generatedId := uuid.New().String()

	_, err := r.db.Exec(query, generatedId, note.Type, note.Body, note.ParticipantId, note.AuthorId)

	if err != nil {
		return "", err
	}

	return generatedId, nil
}
