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

func FetchSublevelsHandler(c echo.Context) error {
	var err = godotenv.Load()
	if err != nil {
		c.Logger().Fatal(err)
		return err
	}
	dbName := os.Getenv("dbName")
	var level_name string
	var cursor *mongo.Cursor


	sublevelsCollection := configs.Client.Database(dbName).Collection(configs.SubLevelsCollection)

	if c.Param("level") != "" {
		level_name = c.Param("level")
		cursor, err = sublevelsCollection.Find(context.Background(), bson.M{"level_name": level_name})
	} else {
		cursor, err = sublevelsCollection.Find(context.Background(), bson.D{})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer cursor.Close(context.Background())

	var subItems []models.Sublevel
	for cursor.Next(context.Background()) {
		var subItem models.Sublevel
		if err := cursor.Decode(&subItem); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		subItems = append(subItems, subItem)
	}

	return c.JSON(http.StatusOK, subItems)
}