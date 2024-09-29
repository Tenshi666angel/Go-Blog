package types

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
    jwt.StandardClaims
    Username string `json:"username"`
}
