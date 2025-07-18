
package repository_test

import (
	"competitions/config"
	"competitions/models"
	"competitions/repository"
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	testRepo repository.TorneioRepository
	testDB   *pgxpool.Pool
)

// TestMain é executado antes de todos os testes neste pacote.
// É usado para configurar o ambiente de teste (banco de dados) e limpar depois.
func TestMain(m *testing.M) {
	// Carrega o .env para ter acesso às variáveis de ambiente do DB
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente existentes")
	}

	// Conecta ao banco de dados de teste
	config.ConnectDatabase()
	testDB = config.DB
	testRepo = repository.NewTorneioRepository(testDB)

	// Executa os testes
	exitCode := m.Run()

	// Limpa e fecha a conexão
	truncateTables()
	testDB.Close()

	os.Exit(exitCode)
}

// truncateTables limpa os dados das tabelas para garantir um estado limpo entre os testes.
func truncateTables() {
	// A ordem é importante por causa das foreign keys
	tables := []string{
		"jogadores_torneios",
		"torneios",
		"esportes",
		"cidades",
		"estados",
		"paises",
	}
	for _, table := range tables {
		_, err := testDB.Exec(context.Background(), "TRUNCATE TABLE "+table+" RESTART IDENTITY CASCADE")
		if err != nil {
			log.Fatalf("Falha ao truncar tabela %s: %v", table, err)
		}
	}
}

// seedDependencies insere dados básicos necessários para criar um torneio.
func seedDependencies(t *testing.T) (int, int, int, int) {
	var paisID, estadoID, cidadeID, esporteID int

	err := testDB.QueryRow(context.Background(), "INSERT INTO paises (nome) VALUES ('Brasil') RETURNING id").Scan(&paisID)
	assert.NoError(t, err)

	err = testDB.QueryRow(context.Background(), "INSERT INTO estados (nome, sigla, id_pais) VALUES ('São Paulo', 'SP', $1) RETURNING id", paisID).Scan(&estadoID)
	assert.NoError(t, err)

	err = testDB.QueryRow(context.Background(), "INSERT INTO cidades (nome, id_estado) VALUES ('São Paulo', $1) RETURNING id", estadoID).Scan(&cidadeID)
	assert.NoError(t, err)

	err = testDB.QueryRow(context.Background(), "INSERT INTO esportes (nome) VALUES ('Tenis') RETURNING id").Scan(&esporteID)
	assert.NoError(t, err)

	return paisID, estadoID, cidadeID, esporteID
}

func TestIntegrationTorneioRepository(t *testing.T) {
	// Garante que a tabela está limpa antes de começar
	truncateTables()
	paisID, estadoID, cidadeID, esporteID := seedDependencies(t)

	ctx := context.Background()
	input := models.TorneioInput{
		Nome:       "Copa do Mundo de Testes",
		DataInicio: time.Now().AddDate(0, 1, 0),
		DataFim:    time.Now().AddDate(0, 2, 0),
		EsporteID:  esporteID,
		CidadeID:   cidadeID,
		EstadoID:   estadoID,
		PaisID:     paisID,
	}

	var createdTorneio models.Torneio

	t.Run("Create", func(t *testing.T) {
		torneio, err := testRepo.Create(ctx, input)

		assert.NoError(t, err)
		assert.NotZero(t, torneio.ID)
		assert.Equal(t, input.Nome, torneio.Nome)
		assert.Equal(t, input.EsporteID, torneio.EsporteID)

		createdTorneio = torneio // Salva para usar nos próximos testes
	})

	t.Run("FindByID", func(t *testing.T) {
		foundTorneio, err := testRepo.FindByID(ctx, createdTorneio.ID)

		assert.NoError(t, err)
		assert.Equal(t, createdTorneio.ID, foundTorneio.ID)
		assert.Equal(t, createdTorneio.Nome, foundTorneio.Nome)
	})

	t.Run("FindAll", func(t *testing.T) {
		torneios, err := testRepo.FindAll(ctx)

		assert.NoError(t, err)
		assert.NotEmpty(t, torneios)
		assert.Len(t, torneios, 1)
		assert.Equal(t, createdTorneio.Nome, torneios[0].Nome)
	})

	t.Run("Update", func(t *testing.T) {
		updateInput := models.TorneioInput{
			Nome:       "Copa do Mundo de Testes - Editado",
			DataInicio: createdTorneio.DataInicio,
			DataFim:    createdTorneio.DataFim,
			EsporteID:  createdTorneio.EsporteID,
			CidadeID:   createdTorneio.CidadeID,
			EstadoID:   createdTorneio.EstadoID,
			PaisID:     createdTorneio.PaisID,
		}

		rowsAffected, err := testRepo.Update(ctx, createdTorneio.ID, updateInput)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)

		// Verifica se o dado foi realmente atualizado
		updatedTorneio, err := testRepo.FindByID(ctx, createdTorneio.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Copa do Mundo de Testes - Editado", updatedTorneio.Nome)
	})

	t.Run("Delete", func(t *testing.T) {
		rowsAffected, err := testRepo.Delete(ctx, createdTorneio.ID)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)

		// Verifica se o torneio foi realmente deletado
		_, err = testRepo.FindByID(ctx, createdTorneio.ID)
		assert.Error(t, err, "Deveria retornar um erro ao buscar torneio deletado")
	})
}
