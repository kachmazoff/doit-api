package mysql

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/model"
)

type ChallengesMysqlRepo struct {
	db *sqlx.DB
}

func NewChallengesMysqlRepo(db *sqlx.DB) *ChallengesMysqlRepo {
	return &ChallengesMysqlRepo{db: db}
}

func (r *ChallengesMysqlRepo) Create(newChallenge model.Challenge) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s (id, author_id, title, body) values (?, ?, ?, ?)", challengesTable)
	generatedId := uuid.New().String()

	_, err := r.db.Exec(
		query,
		generatedId,
		newChallenge.AuthorId,
		newChallenge.Title,
		newChallenge.Body,
	)

	if err != nil {
		return "", err
	}

	return generatedId, nil
}

func (r *ChallengesMysqlRepo) GetAll() ([]model.Challenge, error) {
	query := fmt.Sprintf("SELECT * FROM %s", challengesTable)

	var challenges []model.Challenge
	if err := r.db.Select(&challenges, query); err != nil {
		return []model.Challenge{}, err
	}

	return challenges, nil
}
