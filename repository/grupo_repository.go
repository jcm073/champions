package repository

import (
	"competitions/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GrupoRepository define a interface para interagir com os dados dos grupos.
type GrupoRepository interface {
	CreateGrupos(ctx context.Context, torneioID int, input models.CriarGruposInput) ([]models.GrupoComJogadores, error)
}

// pgGrupoRepository é a implementação concreta para GrupoRepository.
type pgGrupoRepository struct {
	db *pgxpool.Pool
}

// NewGrupoRepository cria uma nova instância de GrupoRepository.
func NewGrupoRepository(db *pgxpool.Pool) GrupoRepository {
	return &pgGrupoRepository{db: db}
}

// CreateGrupos cria grupos para um torneio e categoria, distribuindo os jogadores.
func (r *pgGrupoRepository) CreateGrupos(ctx context.Context, torneioID int, input models.CriarGruposInput) ([]models.GrupoComJogadores, error) {
	// 1. Buscar todos os jogadores inscritos na categoria especificada do torneio.
	queryJogadores := `
		SELECT id_jogador FROM jogadores_torneios
		WHERE id_torneio = $1 AND id_categoria = $2 AND tipo_modalidade = 'simples'
	`
	rows, err := r.db.Query(ctx, queryJogadores, torneioID, input.CategoriaID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar jogadores: %w", err)
	}
	defer rows.Close()

	var playerIDs []int
	for rows.Next() {
		var playerID int
		if err := rows.Scan(&playerID); err != nil {
			return nil, fmt.Errorf("falha ao ler ID do jogador: %w", err)
		}
		playerIDs = append(playerIDs, playerID)
	}

	// 2. Distribuir jogadores em grupos.
	gruposDeJogadores, err := models.DistributePlayers(playerIDs)
	if err != nil {
		return nil, err
	}

	// 3. Iniciar transação.
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback em caso de erro.

	var gruposResult []models.GrupoComJogadores

	// 4. Criar cada grupo e associar jogadores.
	for i, grupoDeJogadores := range gruposDeJogadores {
		// Criar o grupo no banco.
		nomeGrupo := fmt.Sprintf("Grupo %d", i+1)
		var grupoID int
		queryGrupo := `INSERT INTO grupos (id_torneio, id_categoria, nome) VALUES ($1, $2, $3) RETURNING id`
		err := tx.QueryRow(ctx, queryGrupo, torneioID, input.CategoriaID, nomeGrupo).Scan(&grupoID)
		if err != nil {
			return nil, fmt.Errorf("falha ao criar grupo '%s': %w", nomeGrupo, err)
		}

		// Associar jogadores ao grupo.
		var jogadoresDoGrupo []models.Usuario
		for _, jogadorID := range grupoDeJogadores {
			queryAssociacao := `INSERT INTO grupo_jogadores_torneios (id_grupo, id_jogador_torneio) VALUES ($1, $2)`
			_, err := tx.Exec(ctx, queryAssociacao, grupoID, jogadorID)
			if err != nil {
				return nil, fmt.Errorf("falha ao associar jogador %d ao grupo %d: %w", jogadorID, grupoID, err)
			}

			// Buscar detalhes do jogador para a resposta.
			var jogador models.Usuario
			queryJogador := `SELECT id, nome, email FROM usuarios WHERE id = (SELECT id_usuario FROM jogadores WHERE id = $1)`
			err = tx.QueryRow(ctx, queryJogador, jogadorID).Scan(&jogador.ID, &jogador.Nome, &jogador.Email)
			if err != nil {
				return nil, fmt.Errorf("falha ao buscar detalhes do jogador %d: %w", jogadorID, err)
			}
			jogadoresDoGrupo = append(jogadoresDoGrupo, jogador)
		}

		gruposResult = append(gruposResult, models.GrupoComJogadores{
			Grupo: models.Grupo{
				ID:          grupoID,
				TorneioID:   torneioID,
				CategoriaID: input.CategoriaID,
				Nome:        nomeGrupo,
			},
			Jogadores: jogadoresDoGrupo,
		})
	}

	// 5. Commit da transação.
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("falha ao commitar transação: %w", err)
	}

	return gruposResult, nil
}
