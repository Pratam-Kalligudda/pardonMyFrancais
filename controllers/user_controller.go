package controllers

import (
	"context"
	"net/http"

	"example.com/backend/configs"
	"example.com/backend/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserCollection() *mongo.Collection {
	return configs.GetClient().Database("pardon_my_francais").Collection("users")
}
// CreateUser inserts a new user into the database
func CreateUser(u models.User) (*mongo.InsertOneResult, error) {
	return UserCollection().InsertOne(context.Background(), u)
}

func UserProgressCollection() *mongo.Collection {
	return configs.GetClient().Database("pardon_my_francais").Collection("user_progress")
}

func CreateUserProgress(userProgress models.UserProgress) (*mongo.InsertOneResult, error) {
	return UserProgressCollection().InsertOne(context.Background(), userProgress)
}

// FindUserByUsername finds a user by username
func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := UserCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func FindUserByUserId(userId string) (*models.User, error) {
	var user models.User
	err := UserCollection().FindOne(context.Background(), bson.M{"userId": userId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
 


func FetchUserHandler(c echo.Context) error {
	userId := c.Get("userId").(string)
	var user models.User

	filter := bson.M{"userId": userId}
	err := UserCollection().FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	responseUser := ConnectUserToResponse(user)
	return c.JSON(http.StatusOK, responseUser)
}

func ConnectUserToResponse(user models.User)models.UserResponse{
	  responseUser := models.UserResponse{
		Username:        user.Username,
		Email:           user.Email,
		RegistrationDate:user.RegistrationDate,
		Name:            user.Name,
		Bio:             user.Bio,
		Location:        user.Location,
		DoB:             user.DoB,
	  }
	  return responseUser
}