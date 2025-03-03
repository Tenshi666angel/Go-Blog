package user

import (
	"blog/internal/types"
)

type UserRepo interface {
	Create(dto UserDto) error
	GetByUsername(username string) (*UserDto, error)
	CreateAvatar(username, url string) error
}

type UserService interface {
	Register(dto UserDto) error
	Login(dto UserDto) error
	CreateAvatar(username string, fileArgs types.FileArgs) (string, error)
}
