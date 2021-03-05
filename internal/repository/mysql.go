package repository

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMysqlDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName))

	if err != nil {
		return nil, err
	}

	return db, nil
}
