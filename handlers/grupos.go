package handlers

import (
	"competitions/models"
	"competitions/repository"
	"net/http"
	"sort"
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
func (h *GrupoHandler) CreateGrupos(c *gin.Context) {
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

// DefinirVencedoresGrupo é o handler para definir os vencedores de um grupo.
func (h *GrupoHandler) DefinirVencedoresGrupo(c *gin.Context) {
	grupoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do grupo inválido"})
		return
	}

	estatisticas, err := h.repo.GetEstatisticasGrupo(c.Request.Context(), grupoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(estatisticas) < 2 {
		c.JSON(http.StatusOK, gin.H{"message": "Não há jogadores suficientes para definir vencedores"})
		return
	}

	// Ordena os jogadores com base nos critérios
	sort.SliceStable(estatisticas, func(i, j int) bool {
		if estatisticas[i].SetsGanhos != estatisticas[j].SetsGanhos {
			return estatisticas[i].SetsGanhos > estatisticas[j].SetsGanhos
		}
		return estatisticas[i].PontosGanhos > estatisticas[j].PontosGanhos
	})

	// Define os vencedores
	vencedores := make([]models.VencedorGrupo, 2)
	for i := 0; i < 2; i++ {
		playerStats := estatisticas[i]
		criterio := "Total de Sets Ganhos"
		// Verifica se houve empate nos sets e o desempate foi por pontos
		if i > 0 && playerStats.SetsGanhos == estatisticas[i-1].SetsGanhos {
			criterio = "Total de Pontos Conquistados"
		}

		vencedores[i] = models.VencedorGrupo{
			Posicao:      i + 1,
			JogadorID:    playerStats.JogadorID,
			NomeJogador:  playerStats.NomeJogador,
			Criterio:     criterio,
			SetsGanhos:   playerStats.SetsGanhos,
			PontosGanhos: playerStats.PontosGanhos,
		}
	}

	c.JSON(http.StatusOK, models.ResultadoVencedores{Vencedores: vencedores})
}

