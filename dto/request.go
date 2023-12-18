package dto

import "github.com/dgrijalva/jwt-go"

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtClaim struct {
	Email string
	ID    int
	jwt.StandardClaims
}
