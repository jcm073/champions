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
	GetEstatisticasGrupo(ctx context.Context, grupoID int) ([]models.EstatisticasJogador, error)
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
	// 1. Buscar todos os jogadores inscritos na categoria especificada do torneio, ordenados por rating.
	queryJogadores := `
		SELECT jt.id_jogador, s.rating
		FROM jogadores_torneios jt
		JOIN jogadores j ON jt.id_jogador = j.id
		JOIN scouts s ON j.id_scout = s.id
		WHERE jt.id_torneio = $1 AND jt.id_categoria = $2 AND jt.tipo_modalidade = 'simples'
		ORDER BY s.rating DESC
	`
	rows, err := r.db.Query(ctx, queryJogadores, torneioID, input.CategoriaID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar jogadores: %w", err)
	}
	defer rows.Close()

	var rankedPlayers []models.RankedPlayer
	for rows.Next() {
		var p models.RankedPlayer
		if err := rows.Scan(&p.ID, &p.Rating); err != nil {
			return nil, fmt.Errorf("falha ao ler ID e rating do jogador: %w", err)
		}
		rankedPlayers = append(rankedPlayers, p)
	}

	// 2. Distribuir jogadores em grupos.
	gruposDeJogadores, err := models.DistributePlayersRanked(rankedPlayers)
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

func (r *pgGrupoRepository) GetEstatisticasGrupo(ctx context.Context, grupoID int) ([]models.EstatisticasJogador, error) {
	query := `
		SELECT 
			j.id AS jogador_id,
			u.nome AS nome_jogador,
			COALESCE(SUM(CASE WHEN s.vencedor_set = j.id THEN 1 ELSE 0 END), 0) AS sets_ganhos,
			COALESCE(SUM(s.pontos_jogador1 + s.pontos_jogador2), 0) AS pontos_ganhos
		FROM grupo_jogadores_torneios g_jt
		JOIN jogadores_torneios jt ON g_jt.id_jogador_torneio = jt.id
		JOIN jogadores j ON jt.id_jogador = j.id
		JOIN usuarios u ON j.id_usuario = u.id
		LEFT JOIN jogos jg ON jg.id_grupo = g_jt.id_grupo AND (jg.id_jogador_torneio1 = jt.id OR jg.id_jogador_torneio2 = jt.id)
		LEFT JOIN sets s ON s.id_jogo = jg.id
		WHERE g_jt.id_grupo = $1
		GROUP BY j.id, u.nome
		ORDER BY sets_ganhos DESC, pontos_ganhos DESC;
	`

	rows, err := r.db.Query(ctx, query, grupoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar estatísticas do grupo: %w", err)
	}
	defer rows.Close()

	var estatisticas []models.EstatisticasJogador
	for rows.Next() {
		var e models.EstatisticasJogador
		if err := rows.Scan(&e.JogadorID, &e.NomeJogador, &e.SetsGanhos, &e.PontosGanhos); err != nil {
			return nil, fmt.Errorf("falha ao ler estatísticas do jogador: %w", err)
		}
		estatisticas = append(estatisticas, e)
	}

	return estatisticas, nil
}
