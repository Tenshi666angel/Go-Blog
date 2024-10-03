package sqlite

import (
	"blog/internal/types"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func (s *Storage) SaveUser(username string, password string) (int64, error) {
	const op = "persistence.sqlite.Register"

	stmt, err := s.db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(username, password)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok &&
			sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	defer stmt.Close()

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetByUsername(username string) (*types.User, error) {
	const op = "persistence.sqlite.GetByUsername"

	stmt, err := s.db.Prepare("SELECT username, password FROM users WHERE username = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user types.User

	if err := stmt.QueryRow(username).Scan(&user.Username, &user.Password); errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %s not found", username)
	}

	defer stmt.Close()

	return &user, nil
}

func (s *Storage) GetEntityById(id int64) (*types.UserEntity, error) {
	const op = "persistence.sqlite.GetEntityById"

	stmt, err := s.db.Prepare("SELECT id, username, password FROM users WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user types.UserEntity

	if err := stmt.QueryRow(id).Scan(&user.Id, &user.Username, &user.Password); errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %d not found", id)
	}

	defer stmt.Close()

	return &user, nil
}

func (s *Storage) GetEntityByUsername(username string) (*types.UserEntity, error) {
	const op = "persistence.sqlite.GetEntityByUsername"

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE username = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user types.UserEntity

	if err := stmt.QueryRow(username).Scan(&user.Id, &user.Username, &user.Password); errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %s not found: %w", username, err)
	}

	defer stmt.Close()

	return &user, nil
}
