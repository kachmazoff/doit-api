package mysql

import (
	"fmt"

	"github.com/google/uuid"
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

func (r *ParticipantsMysqlRepo) Create(participant model.Participant) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s (id, challenge_id, user_id, anonymous, visible_type) values (?, ?, ?, ?, ?)", participantsTable)
	generatedId := uuid.New().String()

	_, err := r.db.Exec(query, generatedId, participant.ChallengeId, participant.UserId, participant.Anonymous, participant.VisibleType)

	if err != nil {
		return "", err
	}

	return generatedId, nil
}

func (r *ParticipantsMysqlRepo) GetActivePublicParticipantsOfUser(userId string) ([]model.Participant, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? AND status='in_progress' AND anonymous=false AND visible_type='public' ORDER BY created DESC", participantsTable)

	var participants []model.Participant
	if err := r.db.Select(&participants, query); err != nil {
		return []model.Participant{}, err
	}

	return participants, nil
}

func (r *ParticipantsMysqlRepo) GetActivePublicParticipantsInChallenge(challengeId string) ([]model.Participant, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE challenge_id=? AND status='in_progress' AND visible_type='public' ORDER BY created DESC", participantsTable)

	var participants []model.Participant
	if err := r.db.Select(&participants, query); err != nil {
		return []model.Participant{}, err
	}

	return participants, nil
}
