package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kachmazoff/doit-api/internal/repository/mysql"
)

type Repositories struct {
	Users Users
}

func NewMysqlRepos(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users: mysql.NewUsersMysqlRepo(db),
	}
}
