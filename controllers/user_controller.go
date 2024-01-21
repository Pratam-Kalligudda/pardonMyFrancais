package controllers

import (
	"context"
	"net/http"
	"os"

	"example.com/backend/configs"
	"example.com/backend/models"
	"github.com/joho/godotenv"
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

// FindUserByUsername finds a user by username
func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := UserCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func FetchUserHandler(c echo.Context) error {
	var authUser string
	var user models.User
	var err = godotenv.Load()
	if err != nil {
		c.Logger().Fatal(err)
		return err
	}
	dbName := os.Getenv("dbName")
	collection := configs.Client.Database(dbName).Collection(configs.UsersCollection)
	authUser = c.Param("user")
	filter := bson.M{"username": authUser}
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	return c.JSON(http.StatusOK, user)
}