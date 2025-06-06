package handlers

import (
	"competitions/config"
	"competitions/models"
	"competitions/utils"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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

	if err := config.DB.Create(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível criar o usuario"})
		return
	}

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

	if err := config.DB.Where("email = ?", loginData.Email).First(&usuario).Error; err != nil {
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

func Logout(c *gin.Context) {
	// For stateless JWT authentication, logout is handled on the client side.
	// Here we can just return a success message.
	// Logic to handle user logout
	// This could involve invalidating the JWT token or clearing session data
	c.SetCookie("token", "", -1, "/", "localhost", false, true) // Clear the cookie
	c.JSON(http.StatusOK, gin.H{"message": "Logout realizado com sucesso"})
}
