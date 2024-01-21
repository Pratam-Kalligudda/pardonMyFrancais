package models

	import (
		"github.com/dgrijalva/jwt-go"
	)
type JWTToken struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
