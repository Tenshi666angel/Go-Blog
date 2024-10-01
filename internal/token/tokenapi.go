package token

import (
	"blog/internal/types"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const signingKey = "osue3q97tgtg72gq3tgfvv"

func GenerateToken(username string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: username,
	})

	return token.SignedString([]byte(signingKey))
}

func ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &types.TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*types.TokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type TokenClaims")
	}

	return claims.Username, nil
}

func SetToCookie(token string, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "tasty_cookies",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   30 * 24 * 60 * 60,
	}

	http.SetCookie(w, &cookie)
}
