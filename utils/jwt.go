package utils

import (
	"competitions/config"
	"competitions/models"
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
func JwtMiddleware() (*jwt.GinJWTMiddleware, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: .env não encontrado ou não pôde ser carregado")
	}
	secretKey := os.Getenv("JWT_SECRET")
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
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return map[string]interface{}{
				"user_id": claims["user_id"],
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginData struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			if err := c.ShouldBindJSON(&loginData); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			email := loginData.Email
			password := loginData.Password

			query := `
				SELECT id, username, email, password
				FROM usuarios
				WHERE email = $1
			`
			var user models.Usuario
			err := config.DB.QueryRow(context.Background(), query, email).
				Scan(&user.ID, &user.Username, &user.Email, &user.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return map[string]interface{}{
				"user_id": user.ID,
			}, nil
		},
		Authorizator: func(data any, c *gin.Context) bool {
			// Implemente sua lógica de autorização aqui
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
