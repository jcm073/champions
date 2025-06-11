package handlers

import (
	"competitions/config"
	"competitions/models"
	"competitions/utils"
	"context"
	"fmt"
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
// @Param        user  body      models.Usuario  true  "User data"
// @Success      201  {object}  gin.H{"message": "Usuario criado com sucesso", "usuario": models.Usuario, "token": string}
// @Failure      400  {object}  gin.H{"error": "Error message"}
// @Failure      500  {object}  gin.H{"error": "Error message"}
func Signup(c *gin.Context) {
	var usuario models.Usuario

	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ShouldBindJSON": err.Error()})
		return
	}

	fmt.Println("RAW password received:", usuario.Password)

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
	var id int
	err = config.DB.QueryRow(
		context.Background(),
		query,
		usuario.Tipo, usuario.Nome, usuario.Username, usuario.CPF, usuario.DataNascimento,
		usuario.Email, usuario.Password, usuario.Telefone, usuario.Instagram, usuario.Ativo,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível criar o usuario: " + err.Error()})
		return
	}
	usuario.ID = uint(id)

	token, err := utils.GenerateJWT(usuario.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível gerar o token"})
		return
	}
	c.SetCookie("token", token, 24*3600, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario criado com sucesso",
		"usuario": gin.H{
			"id":       usuario.ID,
			"username": usuario.Username,
			"email":    usuario.Email,
		},
		"token": token,
	})
}

// Login handles user login
// It checks the provided email and password against the database
// If valid, it returns a JWT token
// @Summary      Login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginData  body      struct{Email string; Password string}  true  "Login data"
// @Success      200  {object}  gin.H{"token": string, "usuario": models.Usuario}
// @Failure      400  {object}  gin.H{"error": "Error message"}
// @Failure      401  {object}  gin.H{"error": "Invalido email ou senha"}
// @Failure      500  {object}  gin.H{"error": "Error message"}
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var usuario models.Usuario

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		SELECT id, username, email, password
		FROM usuarios
		WHERE email = $1
	`
	err := config.DB.QueryRow(context.Background(), query, loginData.Email).
		Scan(&usuario.ID, &usuario.Username, &usuario.Email, &usuario.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalido email ou senha"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalido email ou senha"})
		return
	}

	token, err := utils.GenerateJWT(usuario.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível gerar o token"})
		return
	}

	c.SetCookie("token", token, 24*3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"usuario": gin.H{
			"id":       usuario.ID,
			"username": usuario.Username,
			"email":    usuario.Email,
		},
	})

}

// Logout handles user logout
// It clears the JWT token cookie
// @Summary      Logout
// @Description  Logout user and clear JWT token
// @Tags         auth
// @Produce      json
// @Success      200  {object}  gin.H{"message": "Logout realizado com sucesso"}
// @Failure      500  {object}  gin.H{"error": "Error message"}
// @Security     ApiKeyAuth
// @Security     BearerAuth
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout realizado com sucesso"})
}
