package handlers

import (
	"competitions/models"
	"competitions/repository"
	"errors"
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo repository.UsuarioRepository
}

func NewAuthHandler(userRepo repository.UsuarioRepository) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo}
}

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "user_id"

// Login é a função autenticadora para o middleware JWT.
func (h *AuthHandler) Login(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	password := loginVals.Password

	user, err := h.UserRepo.FindByEmail(c.Request.Context(), email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, jwt.ErrFailedAuthentication
		}
		log.Printf("Erro ao buscar usuário por e-mail '%s': %v", email, err)
		return nil, jwt.ErrFailedAuthentication
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return user, nil
}

// Payload é a função que define o que vai dentro do token JWT.
func (h *AuthHandler) Payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.Usuario); ok {
		return jwt.MapClaims{
			identityKey: v.ID,
			"type":      v.Tipo,
		}
	}
	return jwt.MapClaims{}
}

// IdentityHandler extrai a identidade do usuário do token.
func (h *AuthHandler) IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &models.Usuario{
		ID:   uint(claims[identityKey].(float64)),
		Tipo: claims["type"].(string),
	}
}

// Authorizator verifica se o usuário tem permissão para acessar a rota.
func (h *AuthHandler) Authorizator(data interface{}, c *gin.Context) bool {
	if _, ok := data.(*models.Usuario); ok {
		return true
	}
	return false
}

// Unauthorized é chamado quando a autenticação falha.
func (h *AuthHandler) Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
