package mysql

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/model"
)

type TimelineMysqlRepo struct {
	db *sqlx.DB
}

func NewTimelineMysqlRepo(db *sqlx.DB) *TimelineMysqlRepo {
	return &TimelineMysqlRepo{db: db}
}

func (r *TimelineMysqlRepo) CreateChallenge() error {
	// query := fmt.Sprintf("INSERT INTO %s () VALUES ()", timelineTable)
	return errors.New("Not implemented")
}

func (r *TimelineMysqlRepo) GetAll() ([]model.TimelineItem, error) {
	query := fmt.Sprintf(`
	%s
	WHERE
		(t.challenge_id IS NOT NULL AND c.visible_type='public') AND (
			(c.show_author = true AND t.participant_id IS NULL AND t.note_id IS NULL AND t.suggestion_id IS NULL) 
			OR
			(t.participant_id IS NOT NULL AND p.visible_type = 'public' AND (
				(t.suggestion_id IS NULL AND p.anonymous = false AND (t.note_id IS NOT NULL OR c.show_author = false OR p.user_id != c.author_id))
				OR
				(t.suggestion_id IS NOT NULL AND t.note_id IS NULL and s.anonymous = false)
			))
		)
	ORDER BY t.created DESC
	`, selectBaseTimelineQuery())

	return r.selectTimeline(query)
}

func (r *TimelineMysqlRepo) GetCommon() ([]model.TimelineItem, error) {
	query := fmt.Sprintf(`
	%s
	WHERE
		(t.challenge_id IS NOT NULL AND c.visible_type='public') AND (
			(t.participant_id IS NULL AND t.note_id IS NULL AND t.suggestion_id IS NULL) 
			OR
			(t.participant_id IS NOT NULL AND p.visible_type = 'public' AND (
				(t.suggestion_id IS NULL AND (t.note_id IS NOT NULL OR c.show_author = false OR p.user_id != c.author_id))
				OR
				(t.suggestion_id IS NOT NULL AND t.note_id IS NULL)
			))
		)
	ORDER BY t.created DESC
	`, selectBaseTimelineQuery())

	return r.selectTimeline(query)
}

func (r *TimelineMysqlRepo) GetForUser(userId string) ([]model.TimelineItem, error) {
	query := fmt.Sprintf(`
	%s
	WHERE
		u.id IN (%s) AND
		t.type='CREATE_CHALLENGE'	AND c.visible_type='public' AND c.show_author=true
		OR
		t.type='ACCEPT_CHALLENGE'	AND c.visible_type='public' AND p.visible_type='public' AND p.anonymous=false AND (c.show_author = false OR p.user_id != c.author_id)
		OR
		t.type='ADD_NOTE'			AND c.visible_type='public' AND p.visible_type='public' AND p.anonymous=false
		OR
		t.type='ADD_SUGGESTION'		AND c.visible_type='public' AND p.visible_type='public' AND s.anonymous=false
	ORDER BY t.created DESC
	`, selectBaseTimelineQuery(), selectFolloweesQuery(userId))

	return r.selectTimeline(query)
}

/*
(t.challenge_id IS NOT NULL AND c.visible_type='public') AND (
	(t.participant_id IS NULL AND t.note_id IS NULL AND t.suggestion_id IS NULL)
	OR
	(t.participant_id IS NOT NULL AND p.visible_type = 'public' AND (
		(t.suggestion_id IS NULL AND (t.note_id IS NOT NULL OR c.show_author = false OR p.user_id != c.author_id))
		OR
		(t.suggestion_id IS NOT NULL AND t.note_id IS NULL)
	))
)
*/

func (r *TimelineMysqlRepo) GetUserOwn(userId string) ([]model.TimelineItem, error) {
	query := fmt.Sprintf(`
	%s
	WHERE
		u.id=? AND (t.type!='ACCEPT_CHALLENGE' OR p.user_id != c.author_id)
	ORDER BY t.created DESC
	`, selectBaseTimelineQuery())

	return r.selectTimeline(query, userId)
}

func (r *TimelineMysqlRepo) selectTimeline(query string, args ...interface{}) ([]model.TimelineItem, error) {
	var timeline []model.TimelineItem
	if err := r.db.Select(&timeline, query, args...); err != nil {
		println(err.Error())
		return []model.TimelineItem{}, err
	}

	if timeline == nil {
		timeline = []model.TimelineItem{}
	}

	return timeline, nil
}

func selectBaseTimelineQuery() string {
	query := fmt.Sprintf(`
	SELECT 
		t.*,
		
		u.id as "user.id",
		u.username as "user.username",
		u.email as "user.email",
		u.created as "user.created",
		
		c.id AS "challenge.id",
		c.created AS "challenge.created",
		c.author_id AS "challenge.author_id",
		c.show_author AS "challenge.show_author",
		c.title AS "challenge.title",
		c.body AS "challenge.body",
		c.visible_type AS "challenge.visible_type",
		c.participants_type AS "challenge.participants_type"
	FROM %s AS t
		LEFT JOIN %s AS u ON t.user_id=u.id
		LEFT JOIN %s AS c ON t.challenge_id=c.id
		LEFT JOIN %s AS p ON t.participant_id=p.id
		LEFT JOIN %s AS s ON t.suggestion_id=s.id
	`, timelineTable, usersTable, challengesTable, participantsTable, suggestionsTable)
	return query
}

// ` WHERE
// t.type='CREATE_CHALLENGE'	AND c.visible_type='public' AND c.show_author=true
// OR
// t.type='ACCEPT_CHALLENGE'	AND c.visible_type='public' AND p.visible_type='public' AND p.anonymous=false AND (c.show_author = false OR p.user_id != c.author_id)
// OR
// t.type='ADD_NOTE'			AND c.visible_type='public' AND p.visible_type='public' AND p.anonymous=false
// OR
// t.type='ADD_SUGGESTION'		AND c.visible_type='public' AND p.visible_type='public' AND s.anonymous=false
// `
