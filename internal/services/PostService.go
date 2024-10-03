package services

import (
	"blog/internal/constants/servererror"
	"blog/internal/lib/logger/sl"
	"blog/internal/persistence"
	"blog/internal/token"
	"blog/internal/types"
	"fmt"
	"log/slog"
)

type PostService struct {
	Logger   *slog.Logger
	PostRepo persistence.PostsRepo
	UserRepo persistence.UserRepo
}

func NewPosts(logger *slog.Logger,
	postRepo persistence.PostsRepo,
	userRepo persistence.UserRepo) *PostService {
	return &PostService{
		Logger:   logger,
		PostRepo: postRepo,
		UserRepo: userRepo,
	}
}

func (s *PostService) Create(postDto types.PostRequest, accessToken string) (*types.PostResponse, error) {
	const op = "services.PostService.Create"

	username, err := token.ParseToken(accessToken)
	if err != nil {
		s.Logger.Error("failed to parse token", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InvalidCrerdentials)
	}

	user, err := s.UserRepo.GetEntityByUsername(username)
	if err != nil {
		s.Logger.Error("failed to get user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.ResourceNotFound)
	}

	newPost := types.PostEntity{
		Title:   postDto.Title,
		Content: postDto.Content,
		User_id: user.Id,
	}

	id, err := s.PostRepo.SavePost(newPost)
	if err != nil {
		s.Logger.Error("failed to save post", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
	}

	_ = id

	return &types.PostResponse{
		Title:       newPost.Title,
		Content:     newPost.Content,
		Author:      user.Username,
		Likes_count: 0,
	}, nil
}

func (s *PostService) GetAll() (*[]types.PostResponse, error) {
	const op = "services.PostService.GetAll"

	postsEntities, err := s.PostRepo.GetAll()
	if err != nil {
		s.Logger.Error("failed to get posts", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
	}

	var posts []types.PostResponse

	for _, pe := range *postsEntities {
		entity, err := s.UserRepo.GetEntityById(pe.User_id)
		if err != nil {
			s.Logger.Error("failed to get user", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
		}

		post := types.PostResponse{
			Title:       pe.Title,
			Content:     pe.Content,
			Author:      entity.Username,
			Likes_count: pe.Likes_count,
		}
		posts = append(posts, post)
	}

	return &posts, nil
}
