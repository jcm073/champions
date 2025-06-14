package routes

import (
	"competitions/handlers"
	"competitions/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	jwtMiddleware, err := utils.JwtMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.POST("/login", jwtMiddleware.LoginHandler)
	router.POST("/signup", handlers.Signup) // Adicionando a rota de signup

	auth := router.Group("/api")
	auth.Use(jwtMiddleware.MiddlewareFunc())
	{
		auth.GET("/usuarios", handlers.GetUsuarios)
		auth.GET("/usuarios/:id", handlers.GetUsuarioByID)
		auth.PUT("/usuarios/:id", handlers.UpdateUsuario)
		auth.DELETE("/usuarios/:id", handlers.DeleteUsuario)
		// ... outras rotas protegidas
	}

	// Rota de logout
	auth.GET("/logout", jwtMiddleware.LogoutHandler)
	// Rota de refresh token
	auth.GET("/refresh_token", jwtMiddleware.RefreshHandler)
	// Rota de health check
	auth.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "API is running",
		})
	})
}
