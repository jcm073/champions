package main

import (
	"log"

	"github.com/gin-contrib/cors"
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

	router := gin.Default()

	// Configuração do CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // Permitir todas as origens (mudar em produção!)
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	// Registra as rotas
	routes.RegisterRoutes(router)
	// Inicia o servidor
	// Escuta e serve na porta 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
