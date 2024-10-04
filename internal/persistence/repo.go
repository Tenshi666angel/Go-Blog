package persistence

import "blog/internal/types"

type UserRepo interface {
	SaveUser(username string, password string) (int64, error)
	GetByUsername(username string) (*types.User, error)
	GetEntityById(id int64) (*types.UserEntity, error)
	GetEntityByUsername(username string) (*types.UserEntity, error)
}

type PostsRepo interface {
	SavePost(post types.PostEntity) (int64, error)
	GetAll() (*[]types.PostEntity, error)
	UpdateLikes(appId string, like int, username string) (bool, error)
	GetLike(appId string, username string) (*types.Like, error)
	CreateLike(appId string, username string) (int64, error)
	DeleteLike(appId string, username string) error
}
