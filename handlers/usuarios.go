package handlers

import (
	"competitions/models"
	"competitions/repository"
	"errors"
	"log"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// UsuarioHandler holds the repository dependency.
type UsuarioHandler struct {
	repo repository.UsuarioRepository
}

// NewUsuarioHandler creates a new UsuarioHandler with the given repository.
func NewUsuarioHandler(repo repository.UsuarioRepository) *UsuarioHandler {
	return &UsuarioHandler{repo: repo}
}

// Buscar todos os usuários
func (h *UsuarioHandler) GetUsuarios(c *gin.Context) {
	usuarios, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		log.Printf("Erro ao buscar usuários: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os usuários."})
		return
	}
	if len(usuarios) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhum usuário encontrado."})
		return
	}
	// O repositório já omite o campo de senha, então a resposta é segura.
	c.JSON(http.StatusOK, usuarios)
}

// Buscar usuário por ID
func (h *UsuarioHandler) GetUsuarioByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	usuario, err := h.repo.FindByID(c.Request.Context(), idInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		log.Printf("Erro ao buscar usuário por ID %d: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar o usuário."})
		return

	}
	c.JSON(http.StatusOK, usuario)
}

// Criar novo usuário
// NOTA: Esta função é similar à `Signup`, mas pode ser usada por um admin para criar usuários.
// A principal diferença é que esta rota deve ser protegida por autenticação/autorização de admin.
func (h *UsuarioHandler) CreateUsuario(c *gin.Context) {
	var input models.UsuarioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Hash da senha antes de salvar no banco
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao gerar hash da senha: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno."})
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

	err = h.repo.Create(c.Request.Context(), &usuario)
	if err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao criar o usuário."})
		return
	}
	// Limpa o campo Password para não expor a senha
	// Isso é importante para não expor informações sensíveis.
	usuario.Password = ""
	// O campo `Password` já tem `omitempty`, então ele não será serializado.
	// Retorne o usuário criado, mas sem a senha.
	c.JSON(http.StatusCreated, usuario)
}

// Atualizar usuário existente
func (h *UsuarioHandler) UpdateUsuario(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input models.UpdateUsuarioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	rowsAffected, err := h.repo.Update(c.Request.Context(), idInt, &input)
	if err != nil {
		log.Printf("Erro ao atualizar usuário %d: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao atualizar o usuário."})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado para atualizar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuário atualizado com sucesso"})
}

// Deletar usuário
func (h *UsuarioHandler) DeleteUsuario(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	rowsAffected, err := h.repo.Delete(c.Request.Context(), idInt)
	if err != nil {
		log.Printf("Erro ao deletar usuário %d: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao deletar o usuário."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado para deletar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário deletado com sucesso"})
}

// ChangePassword handles changing a user's password
// @Summary      Change User Password
// @Description  Changes the password for a specific user, requiring the old password for verification.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        input body models.ChangePasswordInput true "Old and New Passwords"
// @Success      200  {object}  gin.H{"message": "Senha atualizada com sucesso."}
// @Failure      400  {object}  gin.H{"error": "Error message"}
// @Failure      401  {object}  gin.H{"error": "Error message"}
// @Failure      404  {object}  gin.H{"error": "Error message"}
// @Failure      500  {object}  gin.H{"error": "Error message"}
func (h *UsuarioHandler) ChangePassword(c *gin.Context) {
	// Extrai as claims do token JWT. O middleware já garantiu que o token é válido.
	claims := jwt.ExtractClaims(c)
	// O ID do usuário é armazenado como float64 após a decodificação do JSON do token.
	userIDClaim, ok := claims["user_id"].(float64)
	if !ok {
		// Isso não deve acontecer se o middleware estiver configurado corretamente, mas é uma boa verificação.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuário inválido ou ausente no token."})
		return
	}
	// Converte o ID para o tipo usado no modelo (uint).
	userID := uint(userIDClaim)

	var input models.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// 1. Buscar a senha hash atual do banco de dados
	// Usamos um método de repositório específico que busca a senha.
	// O método FindByID padrão omite a senha por segurança.
	user, err := h.repo.FindByIDForAuth(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado."})
			return
		}
		log.Printf("Erro ao buscar senha do usuário %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno ao verificar a senha."})
		return
	}

	// 2. Comparar a senha antiga fornecida com o hash armazenado
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Senha antiga incorreta."})
			return
		}
		log.Printf("Erro ao comparar senhas para o usuário %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno ao verificar a senha."})
		return
	}

	// 3. Gerar um novo hash para a nova senha
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao gerar hash da nova senha para o usuário %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno ao processar a nova senha."})
		return
	}

	// 4. Atualizar a senha no banco de dados
	rowsAffected, err := h.repo.UpdatePassword(c.Request.Context(), userID, string(newHashedPassword))
	if err != nil {
		log.Printf("Erro ao atualizar senha do usuário %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno ao atualizar a senha."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado para atualizar a senha."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Senha atualizada com sucesso."})
}

// The above code is an implementation of handlers for managing usuarios in a web application using Gin and a PostgreSQL database.
// The functions GetUsuarios, GetUsuarioByID, CreateUsuario, UpdateUsuario, DeleteUsuario, are designed to handle HTTP requests for managing usuarios in your application.
// Note: Ensure you have the necessary imports and that your project structure supports these handlers.
// Also, make sure to handle errors and edge cases as per your application's requirements.
