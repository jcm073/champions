package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"competitions/config"
	"competitions/handlers"
	"competitions/repository"
	"competitions/routes"
)

func main() {
	// Conecta ao banco de dados PostgreSQL
	config.ConnectDatabase()
	if config.DB == nil {
		log.Fatal("Falha ao conectar ao banco de dados")
	}

	// 1. Instanciar Repositórios
	userRepo := repository.NewUsuarioRepository(config.DB)
	torneioRepo := repository.NewTorneioRepository(config.DB)
	esporteRepo := repository.NewEsporteRepository(config.DB)

	// 2. Instanciar Handlers, injetando os repositórios
	authHandler := handlers.NewAuthHandler(userRepo)
	userHandler := handlers.NewUsuarioHandler(userRepo)
	torneioHandler := handlers.NewTorneioHandler(torneioRepo)
	esporteHandler := handlers.NewEsporteHandler(esporteRepo)

	router := gin.Default()

	// Configuração do CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // Permitir todas as origens (mudar em produção!)
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	// Registra as rotas, passando os handlers e repositórios necessários
	routes.RegisterRoutes(router, authHandler, userHandler, torneioHandler, esporteHandler, userRepo)
	// Inicia o servidor
	// Escuta e serve na porta 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
