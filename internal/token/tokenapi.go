package token

import (
	"blog/internal/types"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const signingKey = "osue3q97tgtg72gq3tgfvv"

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
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
