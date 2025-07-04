package routes

import (
	"competitions/handlers"
	"competitions/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes configura todas as rotas da aplicação.
func RegisterRoutes(
	router *gin.Engine,
	userHandler *handlers.UsuarioHandler,
	torneioHandler *handlers.TorneioHandler,
	esporteHandler *handlers.EsporteHandler,
	grupoHandler *handlers.GrupoHandler, // Adicionado
	authHandler *handlers.AuthHandler,
	jwtSecret string,
) {
	// Rota para a documentação Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware de autenticação JWT
	// Passa authHandler para o middleware para que ele possa usar authHandler.Login como Authenticator
	authMiddleware := middleware.AuthMiddleware(jwtSecret, authHandler)

	// Rotas de Autenticação (públicas)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authMiddleware.LoginHandler) // Use o LoginHandler fornecido pelo middleware JWT
	}

	// Rotas de Usuários
	userRoutes := router.Group("/usuarios")
	userRoutes.Use(authMiddleware.MiddlewareFunc()) // Proteger rotas de usuário
	{
		userRoutes.POST("", userHandler.CreateUsuario)
		userRoutes.GET("", userHandler.GetUsuarios)
		userRoutes.GET("/:id", userHandler.GetUsuarioByID)
		userRoutes.PUT("/:id", userHandler.UpdateUsuario)
		userRoutes.DELETE("/:id", userHandler.DeleteUsuario)
		userRoutes.PUT("/:id/change-password", userHandler.ChangePassword) // Ativar rota de mudança de senha
		userRoutes.POST("/:id/associar-esporte", userHandler.AssociateEsporte)
	}

	// Rotas de Torneios
	torneioRoutes := router.Group("/torneios")
	torneioRoutes.Use(authMiddleware.MiddlewareFunc()) // Proteger rotas de torneio
	{
		torneioRoutes.POST("", torneioHandler.CreateTorneio)
		torneioRoutes.GET("", torneioHandler.GetTorneios)
		torneioRoutes.GET("/:id", torneioHandler.GetTorneioByID)
		torneioRoutes.PUT("/:id", torneioHandler.UpdateTorneio)
		torneioRoutes.DELETE("/:id", torneioHandler.DeleteTorneio)
		torneioRoutes.POST("/:id/inscrever", torneioHandler.InscreverJogador)
		torneioRoutes.GET("/:id/inscricoes", torneioHandler.ListarInscricoes) // <-- NOVA ROTA
	torneioRoutes.POST("/:id/grupos", grupoHandler.CriarGrupos)
	}

	// Rotas de Esportes
	esporteRoutes := router.Group("/esportes")
	esporteRoutes.Use(authMiddleware.MiddlewareFunc())
	{
		esporteRoutes.GET("", esporteHandler.GetEsportes)
	}
}
