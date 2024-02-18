package models

	import (
		"github.com/dgrijalva/jwt-go"
	)
type JWTToken struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  User		  `json:"user"`
}
