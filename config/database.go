package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func ConnectDatabase() {
	// Carrega variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: .env não encontrado ou não pôde ser carregado")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("Erro: variáveis de ambiente do banco de dados não estão configuradas")
	}
	// Constrói a string de conexão
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	// Conecta ao banco de dados
	DB, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	// Testa a conexão
	if err = DB.Ping(context.Background()); err != nil {
		log.Fatalf("Não foi possível conectar ao PostgreSQL: %v", err)
	}
	log.Println("Conexão com o banco de dados estabelecida com sucesso")
}
