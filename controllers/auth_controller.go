package controllers

import (
	"net/http"
	"time"

	"example.com/backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key")

// GenerateJWT generates a JWT token
func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTToken{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	})

	return token.SignedString(jwtSecret)
}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: jwtSecret})
}

// HashPassword hashes a password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// ComparePasswords compares a hashed password with a plain text password
func ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func SignUp(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	// Check if the username already exists
	_, err := FindUserByUsername(u.Username)
	if err == nil {
		return c.JSON(400, map[string]string{"error": "Username already exists"})
	}

	// Hash the password
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return err
	}

	// Create a new user
	newUser := models.New(u.Username, u.Email, hashedPassword)
	_, err = CreateUser(newUser)
	if err != nil {
		return err
	}

	token, err := GenerateJWT(u.Username)
	if err != nil {
		return err
	}

	return c.JSON(201, models.AuthResponse{Token: token, User: models.User{Username: newUser.Username, Email: newUser.Email, RegistrationDate: newUser.RegistrationDate}})
}

// Login handles user login
func Login(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	// Find the user in the database
	existingUser, err := FindUserByUsername(u.Username)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	// Compare passwords
	err = ComparePasswords(existingUser.Password, u.Password)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	// Create a JWT token
	token, err := GenerateJWT(u.Username)
	if err != nil {
		return err
	}

	return c.JSON(200, models.AuthResponse{Token: token, User: models.User{Username: existingUser.Username, Email: existingUser.Email, RegistrationDate: existingUser.RegistrationDate}})
}

func UpdateUser(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	var updateData models.UpdateUser
	if err := c.Bind(&updateData); err != nil {
		return echo.ErrBadRequest
	}

	// Get existing user data (from database or elsewhere)
	user, err := FindUserByUsername(u.Username)
	if err != nil {
		return err
	}

	// Update user struct with provided fields
	user.Name = updateData.Name
	user.Bio = updateData.Bio
	user.Location = updateData.Location
	user.DoB = updateData.DoB

	// Perform logic to update user information in your database or elsewhere

	return c.JSON(http.StatusOK, user) // Or appropriate response
}
