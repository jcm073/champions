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

func setupUsuarioRouter(handler *UsuarioHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/usuarios", handler.GetUsuarios)
	return r
}

func TestGetUsuariosHandler(t *testing.T) {
	repo := &MockUsuarioRepository{}
	handler := NewUsuarioHandler(repo)
	router := setupUsuarioRouter(handler)

	req, _ := http.NewRequest("GET", "/usuarios", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var usuarios []models.Usuario
	err := json.Unmarshal(w.Body.Bytes(), &usuarios)
	assert.NoError(t, err)
	assert.Equal(t, "Teste", usuarios[0].Nome)
}
