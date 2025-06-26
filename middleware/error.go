package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// AppError representa um erro customizado da aplicação com um código HTTP.
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

// ErrorHandler é um middleware que trata os erros retornados pelos handlers.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Processa a requisição

		// Pega o último erro que ocorreu
		err := c.Errors.Last()
		if err == nil {
			return
		}

		// Verifica se é um erro customizado da nossa aplicação
		var appErr *AppError
		if errors.As(err.Err, &appErr) {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}

		// Trata erros específicos do banco de dados
		if errors.Is(err.Err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Recurso não encontrado."})
			return
		}

		// Para todos os outros erros, retorna um erro 500 genérico
		log.Printf("Erro interno não tratado: %v", err.Err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno no servidor."})
	}
}
