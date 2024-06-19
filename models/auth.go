package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)
type JWTToken struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

type AuthResponse struct {
	Token 	string    `json:"token"`
	User	UserResponse `json:"userResponse"`
}

type UserResponse struct {
	// UserId           string    `json:"userId,omitempty" `
	Username         string    `json:"username,omitempty" `
	Email            string    `json:"email,omitempty" `
	// Password         string    `json:"password,omitempty"`
	RegistrationDate time.Time `json:"registration_date,omitempty" `
	Name             string    `json:"name,omitempty"`
	Bio              string    `json:"bio,omitempty"`
	Location         string    `json:"location,omitempty"`
	DoB              time.Time `json:"dob,omitempty" ` // Use time.Time for date of birth
	ProfilePhoto 	 float32   `json:"profile_photo,omitempty"`
}