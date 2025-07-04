package handlers

import (
	"competitions/models"
	"competitions/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GrupoHandler encapsula a lógica para as rotas de grupos.
type GrupoHandler struct {
	repo repository.GrupoRepository
}

// NewGrupoHandler cria uma nova instância de GrupoHandler.
func NewGrupoHandler(repo repository.GrupoRepository) *GrupoHandler {
	return &GrupoHandler{repo: repo}
}

// CriarGrupos é o handler para a criação de grupos em um torneio.
func (h *GrupoHandler) CriarGrupos(c *gin.Context) {
	torneioID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do torneio inválido"})
		return
	}

	var input models.CriarGruposInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grupos, err := h.repo.CreateGrupos(c.Request.Context(), torneioID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, grupos)
}
