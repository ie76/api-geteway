package database

import (
	"assignment/config"
	"assignment/internal/errors"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var Db *sql.DB

const (
	Test = "test"
	Main = "main"
)

func Init() *errors.Error {
	cfg := config.GetConfig()

	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	Db, err = sql.Open("postgres", connString)
	if err != nil {
		return errors.New(errors.ErrDatabaseError)
	}

	return nil
}

func GetDB() *sql.DB {
	return Db
}
