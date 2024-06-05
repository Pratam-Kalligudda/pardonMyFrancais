package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"time"

	"example.com/backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"

	// "github.com/labstack/echo/v4/middleware"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key")

func GenerateUserID() string {
	return uuid.New().String()
}

// GenerateJWT generates a JWT token
func GenerateJWT(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTToken{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	})

	return token.SignedString(jwtSecret)
}

// func JWTMiddleware() echo.MiddlewareFunc {
// 	return middleware.JWTWithConfig(middleware.JWTConfig{
// 		SigningKey: jwtSecret})
// }

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.ErrUnauthorized
			}

			// Extract token from Authorization header (assuming Bearer scheme)
			tokenString := strings.SplitN(authHeader, " ", 2)[1]

			// Parse the token with claims
			claims := &models.JWTToken{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					return echo.ErrUnauthorized
				}
				return echo.NewHTTPError(http.StatusBadRequest, "invalid token")
			}

			if !token.Valid {
				return echo.ErrUnauthorized
			}

			// Set user ID in context for access by handlers
			c.Set("userId", claims.UserId)

			return next(c)
		}
	}
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
	u.UserId = GenerateUserID()
	// Create a new user
	newUser := models.New(u.UserId, u.Username, u.Email, hashedPassword)
	_, err = CreateUser(newUser)
	if err != nil {
		return err
	}

	userProgress := models.UserProgress{
		UserID:         u.UserId,
		LevelProgress:  models.LevelProgress{CurrentLevel: 1, LevelScores: map[int]float32{}},
		TotalTimeSpent: 0,
		Streak:         0,
		PointsEarned:   0,
		Achievements:   []string{},
		CurrentCombo:   0,
		HighestCombo:   0,
		LastLessonDate: time.Time{},
	}
	println(userProgress.Achievements, userProgress.CurrentCombo, userProgress.HighestCombo, userProgress.LevelProgress.CurrentLevel, userProgress.LevelProgress.LevelScores, userProgress.PointsEarned, userProgress.Streak, userProgress.UserID)
	_, err = CreateUserProgress(userProgress)
	if err != nil {
		return err
	}
	token, err := GenerateJWT(u.UserId)
	if err != nil {
		return err
	}
	responseUser := ConnectUserToResponse(newUser)

	return c.JSON(201, models.AuthResponse{Token: token, User: responseUser})
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
	u.UserId = existingUser.UserId
	// Create a JWT token

	token, err := GenerateJWT(u.UserId)
	if err != nil {
		return err
	}
	responseUser := ConnectUserToResponse(*existingUser)

	return c.JSON(200, models.AuthResponse{Token: token, User: responseUser})
}

func UpdateUserFields(c echo.Context) error {
	// Get userId and field-value pairs from request body
	userId := c.Get("userId").(string)
	var updateReq struct {
		Updates []struct {
			Field    string      `json:"field" binding:"required"`
			Value    interface{} `json:"value" binding:"required"`
			OldValue interface{} `json:"old_value,omitempty"`
		} `json:"updates" binding:"required"`
	}
	if err := c.Bind(&updateReq); err != nil {
		println(2)
		return err // Handle bad request body format
	}

	// Create update documents for each field-value pair
	updateDoc := bson.M{}
	for _, update := range updateReq.Updates {
		if update.Field == "password" {
			var user struct {
				Password string `bson:"password"`
			}
			err := UserCollection().FindOne(context.Background(), bson.M{"userId": userId}).Decode(&user)
			if err != nil {
				return echo.ErrBadRequest // Handle user not found or other errors
			}

			// Assuming oldValue contains the old password
			oldPassword := update.OldValue.(string)
			println(oldPassword)
			println(user.Password)
			err = ComparePasswords(user.Password,oldPassword )
			if err != nil {
				return c.JSON(401, map[string]string{"error": "Invalid credentials"})
			}
			println(update.Value.(string))

			hashedPassword, err := HashPassword(update.Value.(string))
			if err != nil {
				return err
			}
			update.Value = hashedPassword
			println(update.Value.(string))
			println(update.Field)

		}
		if update.Field == "dob" && reflect.TypeOf(update.Value).Kind() == reflect.String {
			dobString := update.Value.(string)
			dob, err := time.Parse("2006-01-02", dobString)
			if err != nil {
				return echo.ErrBadRequest // Handle parsing error (e.g., invalid date format)
			}
			update.Value = dob
		}
		// updates = append(updates, bson.M{
		// 	"$set": bson.M{
		// 		update.Field: update.Value,
		// 	},
		// })
		// println(&updates)
		updateDoc[update.Field] = update.Value
	}
	update := bson.M{"$set": updateDoc}
	updatesJSON, _ := json.Marshal(update)
	println(string(updatesJSON))
	

	// Perform update in MongoDB (replace with your actual update function)
	filter := bson.M{"userId": bson.M{"$eq": userId}} // Filter by userId
	_, err := UserCollection().UpdateMany(context.Background(), filter, update)

	if err != nil {
		// Handle update error (e.g., documents not found, database error)
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Updated Successfully"}) // No content response on success
}

func DeleteUser(c echo.Context) error {
	// Get userId and field-value pairs from request body
	userId := c.Get("userId").(string)

	// Perform update in MongoDB (replace with your actual update function)
	filter := bson.M{"userId": bson.M{"$eq": userId}} // Filter by userId
	print(filter)
	_, err := UserCollection().DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Deleted Successfully"}) // No content response on success
}
