package handlers

import (
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

// Login
type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Login godoc
//	@Summary		Autentica um usuário
//	@Description	Autentica um usuário com e-mail e senha, retornando um token JWT em caso de sucesso.
//	@Tags			Autenticação
//	@Accept			json
//	@Produce		json
//	@Param			input	body		login	true	"Credenciais de Login"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
// Login é a função autenticadora para o middleware JWT.
func (h *AuthHandler) Login(c *gin.Context) (any, error) {
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
