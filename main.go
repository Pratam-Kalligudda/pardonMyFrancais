package main

import (
	"example.com/backend/configs"
	"example.com/backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Connect to MongoDB
	err := configs.ConnectDB()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Routes
	routes.SetRoutes(e)
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
// comment added 