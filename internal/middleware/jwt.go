package middleware

import (
	"blog/pkg/token"
	"net/http"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("tasty_cookies")
		if err != nil {
			http.Error(w, "empty token", http.StatusForbidden)
			return
		}

		if _, err := token.Validate(accessToken.Value); err != nil {
			http.Error(w, "bad token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
