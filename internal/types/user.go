package types

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserEntity struct {
	Id       int64
	Username string
	Password string
}
