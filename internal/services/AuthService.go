package services

import (
	"blog/internal/constants/servererror"
	"blog/internal/lib/logger/sl"
	"blog/internal/persistence"
	"blog/internal/token"
	"blog/internal/types"
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

func (s *AuthService) Register(user types.User) *types.User {
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("failed to hash password", sl.Err(err))
		panic(servererror.InternalError)
	}

	id, err := s.Repo.SaveUser(user.Username, string(passHash))
	if err != nil {
		s.Logger.Error("failed to save user", sl.Err(err))
		panic(servererror.InternalError)
	}

	_ = id

	s.Logger.Info("registered user", slog.String("user", user.Username))

	return &types.User{
		Username: user.Username,
		Password: string(passHash),
	}
}

func (s *AuthService) LogIn(user types.User) string {
	dbUser, err := s.Repo.GetByUsername(user.Username)
	if err != nil {
		s.Logger.Error("user not found", sl.Err(err))
		panic(servererror.InvalidCrerdentials)
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if hashErr != nil {
		s.Logger.Error("invalid password")
		panic(servererror.InvalidCrerdentials)
	}

	token, err := token.GenerateToken(dbUser.Username)
	if err != nil {
		s.Logger.Error("failed to generate token", sl.Err(err))
		panic(servererror.InternalError)
	}

	return token
}
