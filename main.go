package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"competitions/config"
)

func main() {
	// Conecta ao banco de dados MySQL
	config.ConnectDatabase()
	if config.DB == nil {
		log.Fatal("Falha ao conectar ao banco de dados")
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	r.Run() // roda na porta 8080 por padrão
}
