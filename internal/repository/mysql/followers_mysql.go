package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/model"
)

type FollowersMysqlRepo struct {
	db *sqlx.DB
}

func NewFollowersMysqlRepo(db *sqlx.DB) *FollowersMysqlRepo {
	return &FollowersMysqlRepo{db: db}
}

func (r *FollowersMysqlRepo) GetFollowersIds(userId string) ([]string, error) {
	query := selectFollowersQuery(userId)

	var ids []string
	if err := r.db.Select(&ids, query); err != nil {
		return []string{}, err
	}

	return ids, nil
}

func (r *FollowersMysqlRepo) GetFollowedIds(userId string) ([]string, error) {
	query := selectFolloweesQuery(userId)

	var ids []string
	if err := r.db.Select(&ids, query); err != nil {
		return []string{}, err
	}

	return ids, nil
}

func (r *FollowersMysqlRepo) GetFollowers(userId string) ([]model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", usersTable, selectFollowersQuery(userId))

	var users []model.User
	if err := r.db.Select(&users, query); err != nil {
		return []model.User{}, err
	}

	return users, nil
}

func (r *FollowersMysqlRepo) GetFollowees(userId string) ([]model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", usersTable, selectFolloweesQuery(userId))

	var users []model.User
	if err := r.db.Select(&users, query); err != nil {
		return []model.User{}, err
	}

	return users, nil
}

func (r *FollowersMysqlRepo) Subscribe(fromId, toId string) error {
	query := fmt.Sprintf("INSERT INTO %s (follower_id, followee_id) values (?, ?)", followersTable)

	_, err := r.db.Exec(query, fromId, toId)

	return err
}

func (r *FollowersMysqlRepo) Unsubscribe(fromId, toId string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE follower_id=? AND followee_id=?", followersTable)

	_, err := r.db.Exec(query, fromId, toId)

	return err
}

func selectFolloweesQuery(userId string) string {
	return fmt.Sprintf("SELECT followee_id FROM %s WHERE follower_id='%s'", followersTable, userId)
}

func selectFollowersQuery(userId string) string {
	return fmt.Sprintf("SELECT follower_id FROM %s WHERE followee_id='%s'", followersTable, userId)
}
