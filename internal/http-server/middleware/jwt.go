package middleware

import (
	resp "blog/internal/lib/api/response"
	"blog/internal/token"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			render.JSON(w, r, resp.Error("auth header is empty", http.StatusUnauthorized))
			return
		}
		tokenString := strings.Split(authHeader, " ")[1]
		username, err := token.ParseToken(tokenString)
		if err != nil {
			render.JSON(w, r, resp.Error("invalid token", http.StatusUnauthorized))
			return
		}
		_ = username
		next.ServeHTTP(w, r)
	})
}
