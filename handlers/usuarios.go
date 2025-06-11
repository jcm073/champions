package handlers

import (
	"competitions/config"
	"competitions/models"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Buscar todos os usuários
func GetUsuarios(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), "SELECT id, tipo, nome, username, cpf, data_nascimento, email, telefone, instagram, criado_em, ativo FROM usuarios")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários: " + err.Error()})
		return
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		err := rows.Scan(&u.ID, &u.Tipo, &u.Nome, &u.Username, &u.CPF, &u.DataNascimento, &u.Email, &u.Telefone, &u.Instagram, &u.CriadoEm, &u.Ativo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler usuário: " + err.Error()})
			return
		}
		usuarios = append(usuarios, u)
	}
	c.JSON(http.StatusOK, usuarios)
}

// Buscar usuário por ID
func GetUsuarioByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var u models.Usuario
	query := "SELECT id, tipo, nome, username, cpf, data_nascimento, email, telefone, instagram, criado_em, ativo FROM usuarios WHERE id=$1"
	err = config.DB.QueryRow(context.Background(), query, idInt).Scan(&u.ID, &u.Tipo, &u.Nome, &u.Username, &u.CPF, &u.DataNascimento, &u.Email, &u.Telefone, &u.Instagram, &u.CriadoEm, &u.Ativo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// Criar novo usuário
func CreateUsuario(c *gin.Context) {
	var usuario models.Usuario
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := `
		INSERT INTO usuarios (tipo, nome, username, cpf, data_nascimento, email, password, telefone, instagram, ativo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	err := config.DB.QueryRow(
		context.Background(),
		query,
		usuario.Tipo, usuario.Nome, usuario.Username, usuario.CPF, usuario.DataNascimento,
		usuario.Email, usuario.Password, usuario.Telefone, usuario.Instagram, usuario.Ativo,
	).Scan(&usuario.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, usuario)
}

// Atualizar usuário existente
func UpdateUsuario(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var usuario models.Usuario
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := `
		UPDATE usuarios SET tipo=$1, nome=$2, username=$3, cpf=$4, data_nascimento=$5, email=$6, telefone=$7, instagram=$8, ativo=$9
		WHERE id=$10
	`
	_, err = config.DB.Exec(
		context.Background(),
		query,
		usuario.Tipo, usuario.Nome, usuario.Username, usuario.CPF, usuario.DataNascimento,
		usuario.Email, usuario.Telefone, usuario.Instagram, usuario.Ativo, idInt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuário: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuário atualizado com sucesso"})
}

// Deletar usuário
func DeleteUsuario(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	_, err = config.DB.Exec(context.Background(), "DELETE FROM usuarios WHERE id=$1", idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar usuário: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuário deletado com sucesso"})
}

// The above code is an implementation of handlers for managing usuarios in a web application using Gin and a PostgreSQL database.
// The functions GetUsuarios, GetUsuarioByID, CreateUsuario, UpdateUsuario, DeleteUsuario, are designed to handle HTTP requests for managing usuarios in your application.
// Note: Ensure you have the necessary imports and that your project structure supports these handlers.
// Also, make sure to handle errors and edge cases as per your application's requirements.
