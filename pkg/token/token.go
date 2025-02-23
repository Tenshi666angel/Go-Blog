package token

import (
	"blog/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(config.Cfg.JwtKey)

type Pair struct {
	AccessToken  string
	RefreshToken string
}

type claims struct {
	Username string
	jwt.RegisteredClaims
}

func Generate(username string, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		Username: username, 
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func Validate(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {	
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims.Username, nil
	} 

	return "", fmt.Errorf("invalid token")
}

func GeneratePair(username string) (*Pair, error) {
	accessToken, err := Generate(username, time.Minute*5)
	if err != nil {
		return nil, err
	}

	refreshToken, err := Generate(username, time.Hour*24*30)
	if err != nil {
		return nil, err
	}

	return &Pair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
