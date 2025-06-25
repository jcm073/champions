package handlers

import (
	"competitions/models"
	"competitions/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler holds dependencies for authentication handlers.
type AuthHandler struct {
	userRepo repository.UsuarioRepository
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(userRepo repository.UsuarioRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// Signup handles user registration
// It hashes the password and stores the user in the database
// It returns a JWT token upon successful registration
// @Summary      Signup
// @Description  Register a new user with hashed password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        usuario body models.Usuario true "Dados do usuário"
// @Success      201  {object}  gin.H{"message": "Usuario criado com sucesso", "usuario": models.Usuario, "token": string}
// @Failure      400  {object}  gin.H{"error": "Error message"}
// @Failure      500  {object}  gin.H{"error": "Error message"}
func (h *AuthHandler) Signup(c *gin.Context) {
	var input models.UsuarioInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ShouldBindJSON": err.Error()})
		return
	}

	// Verificar se a senha foi enviada
	if input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password é obrigatório"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível gerar o hash da senha"})
		return
	}

	usuario := models.Usuario{
		Tipo:           input.Tipo,
		Nome:           input.Nome,
		Username:       input.Username,
		CPF:            input.CPF,
		DataNascimento: input.DataNascimento,
		Email:          input.Email,
		Password:       string(hashedPassword),
		Telefone:       input.Telefone,
		Instagram:      input.Instagram,
		Ativo:          *input.Ativo,
	}

	err = h.userRepo.Create(c.Request.Context(), &usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível criar o usuario: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario criado com sucesso"})
}

// Logout handles user logout
// It clears the JWT token cookie
