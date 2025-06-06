package routes

import (
	"competitions/handlers"
	middlewares "competitions/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/signup", handlers.Signup)
	router.POST("/login", handlers.Login)

	router.GET("/usuarios", handlers.GetUsuarios)
	router.GET("/usuarios/:id", handlers.GetUsuarioByID)

	auth := router.Group("/")
	auth.Use(middlewares.JWTAuthMiddleware())
	{
		auth.POST("/usuarios", handlers.CreateUsuario)
		auth.PUT("/usuarios/:id", handlers.UpdateUsuario)
		auth.DELETE("/usuarios/:id", handlers.DeleteUsuario)
		auth.POST("/logout", handlers.Logout)

	}
}
