package db

import "database/sql"

type Connection interface {
	Exec(query string, args ...any) (sql.Result, error)
}

type DatabaseConnection struct {
	*sql.DB
}

func New(db *sql.DB) Connection {
	return &DatabaseConnection{
		DB: db,
	}
}
