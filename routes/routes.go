package routes

import (
	"competitions/handlers"
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "user_id"

// RegisterRoutes configura todas as rotas da aplicação, incluindo o middleware de autenticação.
func RegisterRoutes(router *gin.Engine, userHandler *handlers.UsuarioHandler, torneioHandler *handlers.TorneioHandler, esporteHandler *handlers.EsporteHandler, authHandler *handlers.AuthHandler, jwtSecret string) {

	// Inicializa o middleware JWT.
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte(jwtSecret),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     authHandler.Payload,
		IdentityHandler: authHandler.IdentityHandler,
		Authenticator:   authHandler.Login,
		Authorizator:    authHandler.Authorizator,
		Unauthorized:    authHandler.Unauthorized,
		// HTTPStatusMessageFunc é usado para customizar as mensagens de erro.
		HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
			switch {
			case errors.Is(e, jwt.ErrMissingLoginValues):
				return "É necessário fornecer e-mail e senha."
			case errors.Is(e, jwt.ErrFailedAuthentication):
				return "E-mail ou senha incorretos."
			}
			return e.Error()
		},
	})

	if err != nil {
		panic("JWT Error: " + err.Error())
	}

	// Rotas públicas
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/signup", userHandler.CreateUsuario)

	// Rotas protegidas
	api := router.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		// Adicione suas rotas protegidas aqui
		api.GET("/usuarios", userHandler.GetUsuarios)
		api.POST("/usuarios/:id/esportes", userHandler.AssociateEsporte)
		api.GET("/torneios", torneioHandler.GetTorneios)
		// ... etc
	}
}
