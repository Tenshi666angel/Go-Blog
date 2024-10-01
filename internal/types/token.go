package types

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

type TokenPair struct {
	AccesToken   string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
