package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"token": "fake-token"})
	})
	return r
}

func TestLoginHandler(t *testing.T) {
	repo := &MockUsuarioRepository{}
	NewAuthHandler(repo)
	router := setupAuthRouter()

	input := map[string]string{
		"email":    "teste@email.com",
		"password": "senha123",
	}
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
