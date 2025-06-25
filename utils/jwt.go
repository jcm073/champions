package utils

import (
	"competitions/repository"
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Retorna o middleware configurado do gin-jwt
func JwtMiddleware(userRepo repository.UsuarioRepository) (*jwt.GinJWTMiddleware, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: .env não encontrado ou não pôde ser carregado")
	}
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable not set. Application cannot start securely.")
	}
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "competitions zone",
		Key:         []byte(secretKey),
		Timeout:     24 * time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: "user_id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(map[string]interface{}); ok {
				return jwt.MapClaims{
					"user_id": v["user_id"],
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) any {
			claims := jwt.ExtractClaims(c)
			return map[string]interface{}{
				"user_id": claims["user_id"],
			}
		},
		Authenticator: func(c *gin.Context) (any, error) {
			var loginData struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			if err := c.ShouldBindJSON(&loginData); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			email := loginData.Email
			password := loginData.Password

			user, err := userRepo.FindByEmail(context.Background(), email)
			if err != nil {
				log.Printf("Erro ao buscar usuário por email '%s': %v", email, err)
				return nil, jwt.ErrFailedAuthentication
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return map[string]any{
				"user_id": user.ID,
			}, nil
		},
		Authorizator: func(data any, c *gin.Context) bool {
			// Implemente sua lógica de autorização aqui
			// Retorne true se o usuário estiver autorizado, caso contrário, retorne false
			if v, ok := data.(map[string]any); ok {
				userID := v["user_id"]
				if userID == nil {
					return false
				}
				// Aqui você pode verificar se o usuário tem permissão para acessar o recurso
				// Por exemplo, você pode consultar o banco de dados para verificar as permissões do usuário
				// Neste exemplo, vamos apenas permitir todos os usuários autenticados
				log.Printf("Usuário autorizado: %v", userID)
			} else {
				return false
			}
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"token":  token,
				"expire": expire.Unix(),
			})
		},
		// Define onde o token será buscado
		TokenLookup:   "header: Authorization, query: token, cookie: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
