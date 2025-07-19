package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"competitions/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mocks já estão disponíveis pois estão no mesmo package

func setupEsporteRouter(handler *EsporteHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/esportes", handler.GetEsportes)
	return r
}

func TestGetEsportesHandler(t *testing.T) {
	repo := &MockEsporteRepository{}
	handler := NewEsporteHandler(repo)
	router := setupEsporteRouter(handler)

	req, _ := http.NewRequest("GET", "/esportes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var esportes []models.Esporte
	err := json.Unmarshal(w.Body.Bytes(), &esportes)
	assert.NoError(t, err)
	assert.Equal(t, "Futebol", esportes[0].Nome)
}
