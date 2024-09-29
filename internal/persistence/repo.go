package persistence

import "blog/internal/types"

type UserRepo interface {
	SaveUser(username string, password string) (int64, error)
	GetByUsername(username string) (*types.User, error)
}
