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
		u.id IN (%s) AND (
			t.type='CREATE_CHALLENGE'	AND c.visible_type='public' AND c.show_author=true
			OR
			t.type='ACCEPT_CHALLENGE'	AND c.visible_type='public' AND p.visible_type='public' AND p.anonymous=false AND (c.show_author = false OR p.user_id != c.author_id)
			OR
			t.type='ADD_NOTE'			AND c.visible_type='public' AND p.visible_type='public' AND p.anonymous=false
			OR
			t.type='ADD_SUGGESTION'		AND c.visible_type='public' AND p.visible_type='public' AND s.anonymous=false
		)
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

func (r *TimelineMysqlRepo) GetWithFilters(filters model.TimelineFilters) ([]model.TimelineItem, error) {
	var conditions []string

	if filters.TimelineType == "subs" {
		// TODO: validate filters.userId. (can sql injection)
		conditions = append(conditions, fmt.Sprintf(`u.id IN (%s)`, selectFolloweesQuery(filters.UserId)))
	} else if filters.UserId != "" {
		// TODO: validate filters.userId. (can sql injection)
		conditions = append(conditions, fmt.Sprintf(`u.id='%s'`, filters.UserId))
	}

	if filters.LastIndex > -1 {
		operator := "<"
		if filters.OrderType == "ASC" {
			operator = ">"
		}
		conditions = append(conditions, fmt.Sprintf(`t.index %s %d`, operator, filters.LastIndex))
	}
	if filters.EventTypes != nil && hasValidTypes(filters.EventTypes) {
		println("WTF")
		conditions = append(conditions, fmt.Sprintf(`t.type IN %s`, eventTypesToString(filters.EventTypes)))
	}
	if filters.ParticipantId != "" {
		conditions = append(conditions, fmt.Sprintf(`t.participant_id ='%s'`, filters.ParticipantId))
	}
	if filters.ChallengeId != "" {
		conditions = append(conditions, fmt.Sprintf(`t.challenge_id ='%s'`, filters.ChallengeId))
	}

	if filters.TimelineType == "subs" || filters.TimelineType == "common" && filters.UserId != "" && (filters.RequestAuthor == "" || filters.RequestAuthor != filters.UserId) {
		conditions = append(conditions, `(
			t.type='CREATE_CHALLENGE'	AND c.visible_type='public' AND c.show_author=true
			OR
			(t.type='ACCEPT_CHALLENGE' 	OR t.type='ADD_NOTE') AND c.visible_type='public' AND p.visible_type='public' AND p.anonymous=false
			OR
			t.type='ADD_SUGGESTION'		AND c.visible_type='public' AND p.visible_type='public' AND s.anonymous=false
		)`)
	} else if filters.TimelineType == "common" && filters.UserId == "" {
		conditions = append(conditions, `(
			t.type='CREATE_CHALLENGE'	AND c.visible_type='public' 
			OR
			(t.type='ACCEPT_CHALLENGE' 	OR t.type='ADD_NOTE' OR t.type='ADD_SUGGESTION')	AND c.visible_type='public' AND p.visible_type='public' 
		)`)
	}

	sqlConditions := conditions[0]

	for i := 1; i < len(conditions); i++ {
		sqlConditions += fmt.Sprintf(" AND %s", conditions[i])
	}

	query := fmt.Sprintf(`
	%s
	WHERE
		%s
	ORDER BY t.created %s
	LIMIT %d
	`, selectBaseTimelineQuery(), sqlConditions, filters.OrderType, filters.Limit)

	return r.selectTimeline(query)
}

func isValidEventType(eventType string) bool {
	return eventType == "CREATE_CHALLENGE" || eventType == "ACCEPT_CHALLENGE" || eventType == "ADD_NOTE" || eventType == "ADD_SUGGESTION"
}

func hasValidTypes(eventTypes []string) bool {
	res := len(eventTypes) > 0
	for i := 0; i < len(eventTypes); i++ {
		res = res && isValidEventType(eventTypes[i])
	}
	return res
}

func eventTypesToString(eventTypes []string) string {
	eventTypesCollection := "("
	counter := 0
	for i := 0; i < len(eventTypes); i++ {
		if isValidEventType(eventTypes[i]) {
			if counter > 0 {
				eventTypesCollection += ","
			}
			eventTypesCollection += fmt.Sprintf(`'%s'`, eventTypes[i])
		}
	}
	eventTypesCollection += ")"

	return eventTypesCollection
}
