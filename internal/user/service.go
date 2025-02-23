package user

import (
	"blog/pkg/logger/sl"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	logger *slog.Logger
	repo   UserRepo
}

func NewService(logger *slog.Logger, repo UserRepo) *userService {
	return &userService{
		logger: logger,
		repo:   repo,
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
