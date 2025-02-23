package user

type UserRepo interface {
	Create(dto UserDto) error
	GetByUsername(username string) (*UserDto, error)
}

type UserService interface {
	Register(dto UserDto) error
	Login(dto UserDto) error
}
