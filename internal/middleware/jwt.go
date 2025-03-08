package middleware

import (
	"blog/pkg/token"
	"context"
	"net/http"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("tasty_cookies")
		if err != nil {
			http.Error(w, "empty token", http.StatusForbidden)
			return
		}

		username, err := token.Validate(accessToken.Value)
		if err != nil {
			http.Error(w, "bad token", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
