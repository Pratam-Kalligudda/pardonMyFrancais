package routes

import (
	"example.com/backend/controllers"
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo) {
	e.POST("/api/signup", controllers.SignUp)
	e.POST("/api/login", controllers.Login)
	e.POST("/api/updateProfile", controllers.UpdateUser, controllers.JWTMiddleware())
	// e.GET("/api/levels", controllers.FetchLevelsHandler, controllers.JWTMiddleware())
	e.GET("/api/guidebook", controllers.FetchGuidebookHandler)
	e.GET("/api/guidebook/:level", controllers.FetchGuidebookHandler, controllers.JWTMiddleware())
	e.GET("/api/sublevels/:level", controllers.FetchSublevelsHandler, controllers.JWTMiddleware())
	e.GET("/api/sublevels", controllers.FetchSublevelsHandler)
	e.GET("/api/user/:user", controllers.FetchUserHandler, controllers.JWTMiddleware())
	e.GET("/api/audio/:fileName", controllers.FetchAudioHandler)
	e.GET("api/upload/",controllers.UploadHandler)
}
