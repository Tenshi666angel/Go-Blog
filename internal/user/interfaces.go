package user

import "mime/multipart"

type UserRepo interface {
	Create(dto UserDto) error
	GetByUsername(username string) (*UserDto, error)
	CreateAvatar(username, url string) error
}

type UserService interface {
	Register(dto UserDto) error
	Login(dto UserDto) error
	CreateAvatar(
		accessToken string, 
		file multipart.File, 
		handler *multipart.FileHeader) (string, error)
}
