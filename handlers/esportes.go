package handlers

import (
	"competitions/models"
	"competitions/repository"
	"competitions/validation"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// EsporteHandler contém a dependência do repositório.
type EsporteHandler struct {
	repo repository.EsporteRepository
}

// NewEsporteHandler cria um novo EsporteHandler.
func NewEsporteHandler(repo repository.EsporteRepository) *EsporteHandler {
	return &EsporteHandler{repo: repo}
}

func (h *EsporteHandler) CreateEsporte(c *gin.Context) {
	var input models.EsporteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	esporte, err := h.repo.Create(c.Request.Context(), input)
	if err != nil {
		log.Printf("Erro ao criar esporte: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao criar o esporte."})
		return
	}

	c.JSON(http.StatusCreated, esporte)
}

func (h *EsporteHandler) GetEsportes(c *gin.Context) {
	esportes, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		log.Printf("Erro ao buscar esportes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os esportes."})
		return
	}

	c.JSON(http.StatusOK, esportes)
}

func (h *EsporteHandler) GetEsporteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	esporte, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Esporte não encontrado"})
			return
		}
		log.Printf("Erro ao buscar esporte por ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar o esporte."})
		return
	}

	c.JSON(http.StatusOK, esporte)
}

func (h *EsporteHandler) UpdateEsporte(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input models.EsporteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	rowsAffected, err := h.repo.Update(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao atualizar o esporte."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Esporte não encontrado para atualizar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Esporte atualizado com sucesso"})
}

func (h *EsporteHandler) DeleteEsporte(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	rowsAffected, err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		log.Printf("Erro ao deletar esporte %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao deletar o esporte."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Esporte não encontrado para deletar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Esporte deletado com sucesso"})
}
