package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	LevelsCollection    = "levels"
	UsersCollection     = "users"
	SubLevelsCollection = "sublevels"
)

var (
	Client *mongo.Client
)

// ConnectDB connects to MongoDB
func ConnectDB() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	MongoURI := os.Getenv("MongoURI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(MongoURI))
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Check the connection
	err = Client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas!")

	return nil
}

// GetClient returns the MongoDB client
func GetClient() *mongo.Client {
	return Client
}








