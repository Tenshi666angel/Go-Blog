package user

import (
	"blog/pkg/db"
	"blog/pkg/logger/sl"
	"blog/pkg/token"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type UserHandler struct {
	logger  *slog.Logger
	service UserService
}

func NewHandler(logger *slog.Logger, service UserService) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: service,
	}
}

type registerResponse struct {
	Status int
	Msg    string
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "user.handler.Register"

	var req UserDto

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Error("failed to decode json", sl.Err(err))
		http.Error(w, "failed to decode json", http.StatusBadRequest)
		return
	}

	if err := h.service.Register(req); err != nil {
		if errors.Is(err, db.UniqueConstraintError) {
			h.logger.Error("user already exists", sl.Err(err))
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}
		h.logger.Error("failed to register")
		http.Error(w, "failed register user", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, registerResponse{
		Status: http.StatusOK,
		Msg:    "successfully registered",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "user.handler.Login"

	var req UserDto

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Error("failed to decode json", sl.Err(err))
        http.Error(w, "failed to decode json", http.StatusBadRequest)
        return
	}

	if err := h.service.Login(req); err != nil {
		h.logger.Error("bad credentials", sl.Err(err))
		http.Error(w, "bad credentials", http.StatusUnauthorized)
		return
	}

	tokenPair, err := token.GeneratePair(req.Username)
	if err != nil {
		h.logger.Error("error with create tokens", sl.Err(err))
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	c1 := &http.Cookie{
		Name:     "tasty_cookies",
		Value:    tokenPair.AccessToken,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, c1)

	c2 := &http.Cookie{
        Name:     "tasty_cookies2",
        Value:    tokenPair.RefreshToken,
        Secure:   false,
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
    }
    http.SetCookie(w, c2)

	render.JSON(w, r, registerResponse{
		Status: http.StatusOK,
		Msg:    "login successfully",
	})
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	const op = "user.handler.Refresh"

	refreshTokenCookie, err := r.Cookie("tasty_cookies2")
	if err != nil {
		h.logger.Error("empty token", sl.Err(err))
		http.Error(w, "empty token", http.StatusForbidden)
		return
	}

	username, err := token.Validate(refreshTokenCookie.Value)
	if err != nil {
		h.logger.Error("invalid token", sl.Err(err))
		http.Error(w, "bad token", http.StatusForbidden)
		return
	}

	tokenPair, err := token.GeneratePair(username)
    if err != nil {
        h.logger.Error("error with create tokens", sl.Err(err))
        http.Error(w, "500 internal server error", http.StatusInternalServerError)
        return
    }

	c1 := &http.Cookie{
        Name:     "tasty_cookies",
        Value:    tokenPair.AccessToken,
        Secure:   false,
        HttpOnly: false,
        SameSite: http.SameSiteStrictMode,
    }
    http.SetCookie(w, c1)

    c2 := &http.Cookie{
        Name:     "tasty_cookies2",
        Value:    tokenPair.RefreshToken,
        Secure:   false,
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
    }
    http.SetCookie(w, c2)

    render.JSON(w, r, registerResponse{
        Status: http.StatusOK,
        Msg:    "refresh successfully",
    })
}
