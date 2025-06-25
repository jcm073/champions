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

// Handler para criar um torneio

// TorneioHandler holds the repository dependency.
type TorneioHandler struct {
	repo repository.TorneioRepository
}

// NewTorneioHandler creates a new TorneioHandler with the given repository.
func NewTorneioHandler(repo repository.TorneioRepository) *TorneioHandler {
	return &TorneioHandler{repo: repo}
}

func (h *TorneioHandler) CreateTorneio(c *gin.Context) {
	var input models.TorneioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validação aprimorada dos dados de entrada
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	torneio, err := h.repo.Create(c.Request.Context(), input)
	if err != nil {
		log.Printf("Erro ao criar torneio: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno ao criar o torneio."})
		return
	}

	c.JSON(http.StatusCreated, torneio)
}

func (h *TorneioHandler) GetTorneios(c *gin.Context) {
	torneios, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		log.Printf("Erro ao buscar torneios: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os torneios."})
		return
	}

	c.JSON(http.StatusOK, torneios)
}

func (h *TorneioHandler) GetTorneioByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	torneio, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Torneio não encontrado"})
			return
		}
		log.Printf("Erro ao buscar torneio por ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar o torneio."})
		return
	}

	c.JSON(http.StatusOK, torneio)
}

func (h *TorneioHandler) UpdateTorneio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input models.TorneioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validação aprimorada dos dados de entrada
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	rowsAffected, err := h.repo.Update(c.Request.Context(), id, input)
	if err != nil {
		log.Printf("Erro ao atualizar torneio %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao atualizar o torneio."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Torneio não encontrado para atualizar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Torneio atualizado com sucesso"})
}

func (h *TorneioHandler) DeleteTorneio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	rowsAffected, err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		log.Printf("Erro ao deletar torneio %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao deletar o torneio."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Torneio não encontrado para deletar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Torneio deletado com sucesso"})
}
