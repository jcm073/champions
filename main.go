package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"competitions/config"
	"competitions/routes"
)

func main() {
	// Conecta ao banco de dados MySQL
	config.ConnectDatabase()
	if config.DB == nil {
		log.Fatal("Falha ao conectar ao banco de dados")
	}

	r := gin.Default()
	// Registra as rotas
	routes.RegisterRoutes(r)
	// Inicia o servidor
	// Escuta e serve na porta 8080
	r.Run(":8080") // roda na porta 8080 por padrão
}
