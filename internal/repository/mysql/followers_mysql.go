package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FollowersMysqlRepo struct {
	db *sqlx.DB
}

func NewFollowersMysqlRepo(db *sqlx.DB) *FollowersMysqlRepo {
	return &FollowersMysqlRepo{db: db}
}

func (r *FollowersMysqlRepo) GetFollowersIds(userId string) ([]string, error) {
	query := fmt.Sprintf("SELECT follower_id FROM %s WHERE followee_id=?", followersTable)

	var ids []string
	if err := r.db.Select(&ids, query, userId); err != nil {
		return []string{}, err
	}

	return ids, nil
}

func (r *FollowersMysqlRepo) GetFollowedIds(userId string) ([]string, error) {
	query := fmt.Sprintf("SELECT followee_id FROM %s WHERE follower_id=?", followersTable)

	var ids []string
	if err := r.db.Select(&ids, query, userId); err != nil {
		return []string{}, err
	}

	return ids, nil
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
