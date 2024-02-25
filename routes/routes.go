package routes

import (
	"example.com/backend/controllers"
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo) {
	e.POST("/api/signup", controllers.Signup)
	e.POST("/api/login", controllers.Login)
	// e.GET("/api/levels", controllers.FetchLevelsHandler, controllers.JWTMiddleware())
	e.GET("/api/guidebook", controllers.FetchGuidebookHandler)
	e.GET("/api/guidebook/:level", controllers.FetchGuidebookHandler, controllers.JWTMiddleware())
	e.GET("/api/sublevels/:level", controllers.FetchSublevelsHandler, controllers.JWTMiddleware())
	e.GET("/api/sublevels", controllers.FetchSublevelsHandler)
	e.GET("/api/user/:user", controllers.FetchUserHandler, controllers.JWTMiddleware())
	e.GET("/api/audio", controllers.FetchAudioHandler, controllers.JWTMiddleware())
}
