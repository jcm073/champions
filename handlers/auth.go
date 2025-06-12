package handlers

import (
	"competitions/config"
	"competitions/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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
func Signup(c *gin.Context) {
	var usuario models.Usuario

	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ShouldBindJSON": err.Error()})
		return
	}

	// Verificar se a senha foi enviada
	if usuario.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password é obrigatório"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuario.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível gerar o hash da senha"})
		return
	}
	usuario.Password = string(hashedPassword)

	query := `
		INSERT INTO usuarios (tipo, nome, username, cpf, data_nascimento, email, password, telefone, instagram, ativo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
		`
	var id int64 // Alterado para int64
	err = config.DB.QueryRow(context.Background(), query,
		usuario.Tipo, usuario.Nome, usuario.Username, usuario.CPF, usuario.DataNascimento,
		usuario.Email, usuario.Password, usuario.Telefone, usuario.Instagram, usuario.Ativo).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível criar o usuario: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario criado com sucesso"})
}

// Logout handles user logout
// It clears the JWT token cookie
