package routes

import (
	"example.com/backend/controllers"
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo) {
	e.POST("/api/signUp", controllers.SignUp)
	e.POST("/api/logIn", controllers.Login)
	e.POST("/api/updateProfile", controllers.UpdateUserFields, controllers.JWTMiddleware())
	e.POST("/api/updateUserProgress",controllers.UpdateUserProgress, controllers.JWTMiddleware())
	
	e.DELETE("/api/deleteUser", controllers.DeleteUser, controllers.JWTMiddleware())
	
	e.GET("/api/getUserProgress",controllers.GetUserProgress, controllers.JWTMiddleware())
	e.GET("/api/guidebook", controllers.FetchGuidebookHandler)
	e.GET("/api/guidebook/:level", controllers.FetchGuidebookHandler, controllers.JWTMiddleware())
	e.GET("/api/sublevels/:level", controllers.FetchSublevelsHandler, controllers.JWTMiddleware())
	e.GET("/api/sublevels", controllers.FetchSublevelsHandler)
	e.GET("/api/user", controllers.FetchUserHandler, controllers.JWTMiddleware())
	e.GET("/api/audio/:fileName", controllers.FetchAudioHandler)
	e.GET("api/upload",controllers.UploadHandlers)
}
