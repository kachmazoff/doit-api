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

// TODO: Выдавать только то, к чему есть доступ? (challenge.visible_type='public')
func (r *ParticipantsMysqlRepo) GetParticipationsOfUser(userId string, onlyPublic bool, onlyActive bool) ([]model.Participant, error) {
	query := fmt.Sprintf(`
	SELECT 
		p.*,

		c.id AS "challenge.id",
		c.created AS "challenge.created",
		c.author_id AS "challenge.author_id",
		c.show_author AS "challenge.show_author",
		c.title AS "challenge.title",
		c.body AS "challenge.body",
		c.visible_type AS "challenge.visible_type",
		c.participants_type AS "challenge.participants_type"
	FROM %s AS p
		LEFT JOIN %s AS c ON p.challenge_id=c.id 
	WHERE p.user_id=?`, participantsTable, challengesTable)

	if onlyActive {
		query += " AND p.status='in_progress'"
	}
	if onlyPublic {
		query += " AND p.anonymous=false AND p.visible_type='public'"
	}
	query += " ORDER BY p.created DESC"

	return r.selectParticipants(query, userId)
}

func (r *ParticipantsMysqlRepo) GetParticipantsInChallenge(challengeId string, onlyPublic, onlyActive bool) ([]model.Participant, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE challenge_id=?", participantsTable)
	if onlyActive {
		query += " AND status='in_progress'"
	}
	if onlyPublic {
		query += " AND visible_type='public'"
	}
	query += " ORDER BY created DESC"

	return r.selectParticipants(query, challengeId)
}

func (r *ParticipantsMysqlRepo) selectParticipants(query string, args ...interface{}) ([]model.Participant, error) {
	var participants []model.Participant
	if err := r.db.Select(&participants, query, args...); err != nil {
		return []model.Participant{}, err
	}

	if participants == nil {
		participants = []model.Participant{}
	}

	return participants, nil
}
