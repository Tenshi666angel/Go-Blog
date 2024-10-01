package services

import (
	"blog/internal/constants/servererror"
	"blog/internal/lib/logger/sl"
	"blog/internal/persistence"
	"blog/internal/token"
	"blog/internal/types"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Logger *slog.Logger
	Repo   persistence.UserRepo
}

func NewAuth(logger *slog.Logger, repo persistence.UserRepo) *AuthService {
	return &AuthService{
		Logger: logger,
		Repo:   repo,
	}
}

func (s *AuthService) Register(user types.User) (*types.User, error) {
	const op = "services.AuthService.Register"

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("failed to hash password", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
	}

	id, err := s.Repo.SaveUser(user.Username, string(passHash))
	if err != nil {
		s.Logger.Error("failed to save user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
	}

	_ = id

	s.Logger.Info("registered user", slog.String("user", user.Username))

	return &types.User{
		Username: user.Username,
		Password: string(passHash),
	}, nil
}

func (s *AuthService) LogIn(user types.User) (*types.TokenPair, error) {
	const op = "services.AuthService.LogIn"

	dbUser, err := s.Repo.GetByUsername(user.Username)
	if err != nil {
		s.Logger.Error("user not found", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InvalidCrerdentials)
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if hashErr != nil {
		s.Logger.Error("invalid password")
		return nil, fmt.Errorf("%s: %w", op, servererror.InvalidCrerdentials)
	}

	tokenPair, err := token.GeneratePair(dbUser.Username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
	}

	return tokenPair, nil
}

func (s *AuthService) Refresh(refreshToken string) (*types.TokenPair, error) {
	const op = "services.AuthService.Refresh"

	username, err := token.ParseToken(refreshToken)
	if err != nil {
		s.Logger.Error("failed to parse token", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, servererror.InvalidCrerdentials)
	}

	tokenPair, err := token.GeneratePair(username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, servererror.InternalError)
	}

	return tokenPair, nil
}
