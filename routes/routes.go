package routes

import (
	"competitions/handlers"
	"competitions/repository"
	"competitions/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UsuarioHandler,
	torneioHandler *handlers.TorneioHandler,
	esporteHandler *handlers.EsporteHandler,
	userRepo repository.UsuarioRepository,
) {
	jwtMiddleware, err := utils.JwtMiddleware(userRepo)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.POST("/login", jwtMiddleware.LoginHandler)
	router.POST("/signup", authHandler.Signup) // Adicionando a rota de signup

	auth := router.Group("/api")
	auth.Use(jwtMiddleware.MiddlewareFunc())
	{
		auth.POST("/usuarios", userHandler.CreateUsuario) // Rota para admin criar usuários
		auth.GET("/usuarios", userHandler.GetUsuarios)
		auth.GET("/usuarios/:id", userHandler.GetUsuarioByID)
		auth.PUT("/usuarios/:id", userHandler.UpdateUsuario)
		auth.DELETE("/usuarios/:id", userHandler.DeleteUsuario)

		// Rota segura para o usuário logado alterar sua própria senha.
		auth.PUT("/me/password", userHandler.ChangePassword)

		// Rotas de Torneio
		auth.POST("/torneios", torneioHandler.CreateTorneio)
		auth.GET("/torneios", torneioHandler.GetTorneios)
		auth.GET("/torneios/:id", torneioHandler.GetTorneioByID)
		auth.PUT("/torneios/:id", torneioHandler.UpdateTorneio)
		auth.DELETE("/torneios/:id", torneioHandler.DeleteTorneio)

		// Rotas de Esportes (CRUD)
		// Geralmente, estas rotas devem ser protegidas e acessíveis apenas por administradores.
		auth.POST("/esportes", esporteHandler.CreateEsporte)
		auth.GET("/esportes", esporteHandler.GetEsportes)
		auth.GET("/esportes/:id", esporteHandler.GetEsporteByID)
		auth.PUT("/esportes/:id", esporteHandler.UpdateEsporte)
		auth.DELETE("/esportes/:id", esporteHandler.DeleteEsporte)
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
