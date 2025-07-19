package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"competitions/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupGrupoRouter(handler *GrupoHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/torneios/:id/grupos", handler.CreateGrupos)
	return r
}

func TestCreateGruposHandler(t *testing.T) {
	repo := &MockGrupoRepository{}
	handler := NewGrupoHandler(repo)
	router := setupGrupoRouter(handler)

	input := models.CriarGruposInput{
		CategoriaID: 1,
	}
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/torneios/1/grupos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var grupos []models.Grupo
	err := json.Unmarshal(w.Body.Bytes(), &grupos)
	assert.NoError(t, err)
	assert.Equal(t, "Grupo A", grupos[0].Nome)
}
