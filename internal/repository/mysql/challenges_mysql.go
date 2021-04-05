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

func (r *ChallengesMysqlRepo) GetById(id string) (model.Challenge, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=? LIMIT 1", challengesTable)

	var challenge model.Challenge
	if err := r.db.Get(&challenge, query, id); err != nil {
		return model.Challenge{}, err
	}

	return challenge, nil
}

func (r *ChallengesMysqlRepo) GetAll() ([]model.Challenge, error) {
	query := fmt.Sprintf("SELECT * FROM %s", challengesTable)

	var challenges []model.Challenge
	if err := r.db.Select(&challenges, query); err != nil {
		return []model.Challenge{}, err
	}

	if challenges == nil {
		challenges = []model.Challenge{}
	}

	return challenges, nil
}

func (r *ChallengesMysqlRepo) GetAllPublic() ([]model.Challenge, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE visible_type='public'", challengesTable)

	var challenges []model.Challenge
	if err := r.db.Select(&challenges, query); err != nil {
		return []model.Challenge{}, err
	}

	if challenges == nil {
		challenges = []model.Challenge{}
	}

	return challenges, nil
}

func (r *ChallengesMysqlRepo) GetAllOwn(userId string) ([]model.Challenge, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE author_id=?", challengesTable)

	var challenges []model.Challenge
	if err := r.db.Select(&challenges, query, userId); err != nil {
		return []model.Challenge{}, err
	}

	if challenges == nil {
		challenges = []model.Challenge{}
	}

	return challenges, nil
}
