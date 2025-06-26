package middleware

import (
	"competitions/models"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "user_id"

// Authenticator define a interface que a lógica de login deve satisfazer.
// Isso quebra o ciclo de importação entre os pacotes middleware e handlers.
type Authenticator interface {
	Login(c *gin.Context) (interface{}, error)
}

// AuthMiddleware cria e configura o middleware de autenticação JWT.
// Ele é projetado para funcionar com seu `authHandler.Login` customizado,
// focando apenas na validação do token para rotas protegidas.
func AuthMiddleware(secretKey string, auth Authenticator) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Authenticator: auth.Login, // Define o Authenticator usando o método Login da interface
		Realm:         "competitions-api",
		Key:           []byte(secretKey),
		Timeout:       time.Hour * 24,
		MaxRefresh:    time.Hour * 24,
		IdentityKey:   identityKey,
		// PayloadFunc é usado pelo handler de login do middleware para criar o token.
		// Como seu login é customizado, esta função serve como um padrão caso você
		// decida usar o gerador de token da biblioteca em outro lugar.
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.Usuario); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"type":      v.Tipo,
				}
			}
			return jwt.MapClaims{}
		},
		// IdentityHandler extrai a identidade do usuário a partir do token.
		// O valor retornado é passado para a função Authorizator.
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c) // Extrai as claims do token
			// Reconstrói um objeto models.Usuario a partir das claims
			userID, ok := claims[identityKey].(float64)
			if !ok {
				return nil // Ou trate o erro apropriadamente
			}
			userType, ok := claims["type"].(string)
			if !ok {
				return nil // Ou trate o erro apropriadamente
			}
			return &models.Usuario{ID: uint(userID), Tipo: userType} // Retorna um *models.Usuario
		},
		// Authorizator é chamado em cada requisição para verificar se o usuário
		// (identificado pelo token) tem permissão para acessar.
		Authorizator: func(data interface{}, c *gin.Context) bool {
			_, ok := data.(*models.Usuario) // Verifica se os dados são um *models.Usuario
			return ok                       // Se for um Usuario, a autorização é bem-sucedida (por enquanto)
		},
		// Unauthorized é a resposta enviada quando a autenticação falha.
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
		TokenLookup: "header: Authorization",
	})

	if err != nil {
		log.Fatalf("Erro fatal no middleware JWT: %v", err)
	}

	return authMiddleware
}
