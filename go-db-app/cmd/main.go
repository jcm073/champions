package main

import (
    "go-db-app/internal/db"
    "go-db-app/internal/service"
    "go-db-app/pkg/utils/logger"

    "github.com/gin-gonic/gin"
)

func main() {
    // Inicializa conexão com banco de dados
    database, err := db.Connect()
    if err != nil {
        logger.Error("Failed to connect to the database: " + err.Error())
        return
    }
    defer database.Close()

    // Inicializa o serviço de usuário
    userService := service.NewUserService(database)

    // Cria o router Gin
    router := gin.Default()

    // Define as rotas
    router.GET("/users", userService.HandleUsers)

    logger.Info("Starting server on :8080")
    if err := router.Run(":8080"); err != nil {
        logger.Error("Failed to start server: " + err.Error())
    }
}