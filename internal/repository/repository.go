package repository

import (
	"database/sql"
	"testTask/internal/logerr"
)

type Repository struct {
	Db  *sql.DB
	Log *logerr.Logerr
}

func NewRepository(Db *sql.DB, Log *logerr.Logerr) *Repository {
	return &Repository{
		Db:  Db,
		Log: Log,
	}
}

func InitDataBase(dsn string) (*sql.DB, error) {
	Db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = Db.Ping()
	if err != nil {
		return nil, err
	}

	return Db, nil
}
