package mysql

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/model"
)

type UsersMysqlRepo struct {
	db *sqlx.DB
}

func NewUsersMysqlRepo(db *sqlx.DB) *UsersMysqlRepo {
	return &UsersMysqlRepo{db: db}
}

func (r *UsersMysqlRepo) Create(newUser model.User) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s (id, username, password, email) values (?, ?, ?, ?)", usersTable)
	generatedId := uuid.New().String()

	_, err := r.db.Exec(query, generatedId, newUser.Username, newUser.Password, newUser.Email)

	if err != nil {
		return "", err
	}

	return generatedId, nil
}

func (r *UsersMysqlRepo) GetByEmail(email string) (model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=? LIMIT 1", usersTable)

	var user model.User
	if err := r.db.Get(&user, query, email); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *UsersMysqlRepo) GetByUsername(username string) (model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=? LIMIT 1", usersTable)

	var user model.User
	if err := r.db.Get(&user, query, username); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *UsersMysqlRepo) SetStatus(id string, status string) error {
	query := fmt.Sprintf("UPDATE %s SET account_status=? WHERE id=?", usersTable)

	_, err := r.db.Exec(query, status, id)
	if err != nil {
		println(err.Error())
		return err
	}

	return nil
}

func (r *UsersMysqlRepo) GetIdByUsername(username string) (string, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=? LIMIT 1", usersTable)

	var id string
	if err := r.db.Get(&id, query, username); err != nil {
		println(err.Error())
		return "", err
	}

	return id, nil
}

func (r *UsersMysqlRepo) GetEmailById(id string) (string, error) {
	query := fmt.Sprintf("SELECT email FROM %s WHERE id=? LIMIT 1", usersTable)

	var email string
	if err := r.db.Get(&email, query, id); err != nil {
		println(err.Error())
		return "", err
	}

	return email, nil
}

func (r *UsersMysqlRepo) GetAll() ([]model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", usersTable)

	var users []model.User
	if err := r.db.Select(&users, query); err != nil {
		println(err.Error())
		return []model.User{}, err
	}

	return users, nil
}
