package handlers

import (
	"bytes"
	"competitions/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mocks importados de mocks_test.go
// Mocks já estão disponíveis pois estão no mesmo package

func setupTorneioRouter(handler *TorneioHandler) *gin.Engine {
	r := gin.Default()

	// Mocks já estão disponíveis pois estão no mesmo package
	r.POST("/torneios", handler.CreateTorneio)
	return r
}

func TestCreateTorneioHandler(t *testing.T) {
	repo := &MockTorneioRepository{}
	handler := NewTorneioHandler(repo)
	router := setupTorneioRouter(handler)

	input := models.TorneioInput{
		Nome:       "Torneio Teste",
		DataInicio: time.Now(),
		DataFim:    time.Now().Add(24 * time.Hour),
		EsporteID:  1,
		CidadeID:   1,
		EstadoID:   1,
		PaisID:     1,
	}
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/torneios", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var torneio models.Torneio
	err := json.Unmarshal(w.Body.Bytes(), &torneio)
	assert.NoError(t, err)
	assert.Equal(t, "Torneio Teste", torneio.Nome)
}
