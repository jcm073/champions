package main

import (
	"log"
	"os"

	server "github.com/ckanthony/gin-mcp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"competitions/config"
	"competitions/handlers"
	"competitions/middleware"
	"competitions/repository"
	"competitions/routes"

	_ "competitions/docs" // Importa os docs gerados pelo swag
)

//	@title			API para Gerenciamento de Campeonatos
//	@version		1.0
//	@description	Esta é uma API para gerenciar torneios incluindo esportes, equipes e usuários.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	João Martins
//	@contact.email	jmartins@email.com

//	@license.name		BSD-3-Clause
//	@license.version	3.0
//	@license.url		http://www.bsd.org/licenses/bsd-3-clause

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	// Conecta ao banco de dados PostgreSQL
	config.ConnectDatabase()
	if config.DB == nil {
		log.Fatal("Falha ao conectar ao banco de dados")
	}
	// Carrega a chave secreta do JWT do ambiente
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Variável de ambiente JWT_SECRET não definida.")
	}

	// 1. Instanciar Repositórios
	userRepo := repository.NewUsuarioRepository(config.DB)
	torneioRepo := repository.NewTorneioRepository(config.DB)
	esporteRepo := repository.NewEsporteRepository(config.DB)

	grupoRepo := repository.NewGrupoRepository(config.DB)

	// 2. Instanciar Handlers, injetando os repositórios
	authHandler := handlers.NewAuthHandler(userRepo)
	userHandler := handlers.NewUsuarioHandler(userRepo)
	torneioHandler := handlers.NewTorneioHandler(torneioRepo)
	esporteHandler := handlers.NewEsporteHandler(esporteRepo)
	grupoHandler := handlers.NewGrupoHandler(grupoRepo) // Adicionado

	router := gin.Default()

	// Configuração do CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // Permitir todas as origens (mudar em produção!)
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Adiciona o middleware de tratamento de erros
	router.Use(middleware.ErrorHandler())
	// Registra as rotas, passando os handlers e repositórios necessários
	routes.RegisterRoutes(router, userHandler, torneioHandler, esporteHandler, grupoHandler, authHandler, jwtSecret)

	//Create and configure the MCP server
	//Provide essential details for the MCP client.
	mcp := server.New(router, &server.Config{
		Name:        "Championship API",
		Description: "An example API automatically exposed via MCP.",
		// BaseURL is crucial! It tells MCP clients where to send requests.
		BaseURL: "http://localhost:8080",
	})

	//Mount the MCP server endpoint
	mcp.Mount("/mcp") // MCP clients will connect here

	// Inicia o servidor
	// Escuta e serve na porta 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
