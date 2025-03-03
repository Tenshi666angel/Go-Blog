package user

import (
	"blog/internal/dbxutils"
	"blog/internal/types"
	"blog/pkg/logger/sl"
	"fmt"
	"log/slog"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	logger       *slog.Logger
	repo		 UserRepo
	dbxCommitter dbxutils.DbxCommiter
}

func NewService(logger *slog.Logger, repo UserRepo, dbxCommitter dbxutils.DbxCommiter) *userService {
	return &userService{
		logger:		  logger,
		repo:		  repo,
		dbxCommitter: dbxCommitter,
	}
}

func (s *userService) Register(dto UserDto) error {
	const op = "user.service.Register"

	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("error hash password", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return s.repo.Create(UserDto{
		Username: dto.Username,
		Password: string(hash),
	})
}

func (s *userService) Login(dto UserDto) error {
	const op = "user.service.Login"

	user, err := s.repo.GetByUsername(dto.Username)
	if err != nil {
		s.logger.Error("user not found", sl.Err(err))
		return fmt.Errorf("%s: %w", op, BadCredentialsError)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		s.logger.Error("password incorrect")
		return fmt.Errorf("%s: %w", op, BadCredentialsError)
	}

	return nil
}

func (s *userService) CreateAvatar(username string, file types.FileArgs) (string, error) {
	const op = "user.service.CreateAvatar"

	linkChan := make(chan types.StrErr, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	
	go func() {
		defer wg.Done()
		link, err := s.dbxCommitter.Upload(file)
		linkChan <- types.StrErr{Str: link, Err: err}
	}()
	wg.Wait()

	linkRes := <- linkChan
	if linkRes.Err != nil {
		s.logger.Error("error upload avatar", sl.Err(linkRes.Err))
		return "", fmt.Errorf("%s: %w", op, linkRes.Err)
	}

	if err := s.repo.CreateAvatar(username, linkRes.Str); err != nil {
		s.logger.Error("error create avatar", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return linkRes.Str, nil	
}
