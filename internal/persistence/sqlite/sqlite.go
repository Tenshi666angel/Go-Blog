package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "persistence.qlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_username ON users(username);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}
	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	return &Storage{db: db}, nil
}
