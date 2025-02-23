package user

import (
	"blog/pkg/db"
	"blog/pkg/logger/sl"
	"database/sql"	
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

type userRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewRepo(db *sql.DB, logger *slog.Logger) *userRepo {
	return &userRepo{
		db:		db,
		logger: logger,
	}
}

func (r *userRepo) Create(dto UserDto) error {
	const op = "user.repo.Create"

	stmt, err := r.db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		r.logger.Error("error with sql statement", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(dto.Username, dto.Password); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok &&
				mysqlErr.Number == db.UniqueConstraintCode {
			r.logger.Error("error: unique constraint", sl.Err(err))
			return db.UniqueConstraintError
		}
		r.logger.Error("error save user", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *userRepo) GetByUsername(username string) (*UserDto, error) {
	const op = "user.repo.GetByUsername"

	stmt, err := r.db.Prepare("SELECT username, password FROM users WHERE username = ?")
	if err != nil {
		r.logger.Error("error with sql statement", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user UserDto

	if err := stmt.QueryRow(username).Scan(&user.Username, &user.Password); err != nil {
		r.logger.Error("error get user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
