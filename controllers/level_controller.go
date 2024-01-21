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

func FetchGuidebookHandler(c echo.Context) error {
	var level_name string
	var cursor *mongo.Cursor
	var err = godotenv.Load()
	if err != nil {
		c.Logger().Fatal(err)
		return err
	}
	dbName := os.Getenv("dbName")

	collection := configs.Client.Database(dbName).Collection(configs.LevelsCollection)

	if c.Param("level") != "" {
		level_name = c.Param("level")
		cursor, err = collection.Find(context.Background(), bson.M{"level_name": level_name})
	} else {
		cursor, err = collection.Find(context.Background(), bson.D{})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch guidebook content"})
	}
	defer cursor.Close(context.Background())

	var guidebookContents []models.GuidebookContent
	for cursor.Next(context.Background()) {
		var guidebookContent models.GuidebookContent
		if err := cursor.Decode(&guidebookContent); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		guidebookContents = append(guidebookContents, guidebookContent)
	}
	if err := cursor.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, guidebookContents)
}