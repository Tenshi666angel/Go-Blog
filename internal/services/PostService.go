package services

import (
	"blog/internal/constants/servererror"
	"blog/internal/lib/logger/sl"
	"blog/internal/persistence"
	"blog/internal/token"
	"blog/internal/types"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
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
		UserId:  user.Id,
		AppID:   uuid.New().String(),
	}

	id, err := s.PostRepo.SavePost(newPost)
	if err != nil {
		s.Logger.Error("failed to save post", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
	}

	_ = id

	return &types.PostResponse{
		Title:      newPost.Title,
		Content:    newPost.Content,
		Author:     user.Username,
		AppID:      newPost.AppID,
		LikesCount: 0,
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
		entity, err := s.UserRepo.GetEntityById(pe.UserId)
		if err != nil {
			s.Logger.Error("failed to get user", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
		}

		post := types.PostResponse{
			Title:      pe.Title,
			Content:    pe.Content,
			Author:     entity.Username,
			LikesCount: pe.LikesCount,
			AppID:      pe.AppID,
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (s *PostService) Like(accessToken string, appId string) (*types.Like, error) {
	const op = "services.PostService.Like"

	username, err := token.ParseToken(accessToken)
	if err != nil {
		s.Logger.Error("failed to parse token", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InvalidCrerdentials)
	}

	if _, err := s.PostRepo.GetLike(appId, username); err == nil {
		s.Logger.Error("post already liked", slog.String("app_id", appId), slog.String("username", username))
		return nil, fmt.Errorf("post already liked: %w", servererror.BadRequest)
	}

	if _, err := s.PostRepo.CreateLike(appId, username); err != nil {
		s.Logger.Error("failed to create like", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
	}

	_, err = s.PostRepo.UpdateLikes(appId, 1, username)
	if err != nil {
		s.Logger.Error("failed to like user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &types.Like{
		PostID:   appId,
		Username: username,
	}, nil
}

func (s *PostService) UnLike(accessToken string, appId string) (*types.Like, error) {
	const op = "services.PostService.UnLike"

	username, err := token.ParseToken(accessToken)
	if err != nil {
		s.Logger.Error("failed to parse token", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InvalidCrerdentials)
	}

	if _, err := s.PostRepo.GetLike(appId, username); err != nil && errors.Is(err, sql.ErrNoRows) {
		s.Logger.Error("post never liked", slog.String("app_id", appId), slog.String("username", username))
		return nil, fmt.Errorf("post never liked: %w", servererror.BadRequest)
	}

	if err := s.PostRepo.DeleteLike(appId, username); err != nil {
		s.Logger.Error("failed to delete like", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
	}

	_, err = s.PostRepo.UpdateLikes(appId, -1, username)
	if err != nil {
		s.Logger.Error("failed to unlike user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &types.Like{
		PostID:   appId,
		Username: username,
	}, nil
}
