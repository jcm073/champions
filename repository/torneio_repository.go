package repository

import (
	"competitions/models"
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TorneioRepository interface {
	Create(ctx context.Context, input models.TorneioInput) (models.Torneio, error)
	FindAll(ctx context.Context) ([]models.Torneio, error)
	FindByID(ctx context.Context, id int) (models.Torneio, error)
	Update(ctx context.Context, id int, input models.TorneioInput) (int64, error)
	Delete(ctx context.Context, id int) (int64, error)
	InscreverJogador(ctx context.Context, inscricao models.JogadorTorneio) (models.JogadorTorneio, error)
	ListarInscricoesPorTorneio(ctx context.Context, torneioID int) ([]models.InscricaoDetalhada, error)
}

// pgTorneioRepository é a implementação concreta para TorneioRepository.
type pgTorneioRepository struct {
	db *pgxpool.Pool
}

// NewTorneioRepository cria uma nova instância de TorneioRepository.
func NewTorneioRepository(db *pgxpool.Pool) TorneioRepository {
	return &pgTorneioRepository{db: db}
}

// Create insere um novo torneio no banco de dados.
func (r *pgTorneioRepository) Create(ctx context.Context, input models.TorneioInput) (models.Torneio, error) {
	var torneio models.Torneio
	query := `
        INSERT INTO torneios (nome, data_inicio, data_fim, id_esporte, id_cidade, id_estado, id_pais)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, nome, data_inicio, data_fim, id_esporte, id_cidade, id_estado, id_pais, criado_em`
	err := r.db.QueryRow(ctx, query,
		input.Nome, input.DataInicio, input.DataFim, input.EsporteID, input.CidadeID, input.EstadoID, input.PaisID,
	).Scan(
		&torneio.ID, &torneio.Nome, &torneio.DataInicio, &torneio.DataFim,
		&torneio.EsporteID, &torneio.CidadeID, &torneio.EstadoID, &torneio.PaisID, &torneio.CriadoEm,
	)
	return torneio, err
}

// FindAll recupera todos os torneios do banco de dados.
func (r *pgTorneioRepository) FindAll(ctx context.Context) ([]models.Torneio, error) {
	query := `
        SELECT id, nome, data_inicio, data_fim, id_esporte, id_cidade, id_estado, id_pais, criado_em
        FROM torneios
        ORDER BY data_inicio DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Torneio])
}

// FindByID recupera um único torneio pelo seu ID.
func (r *pgTorneioRepository) FindByID(ctx context.Context, id int) (models.Torneio, error) {
	query := `
        SELECT id, nome, data_inicio, data_fim, id_esporte, id_cidade, id_estado, id_pais, criado_em
        FROM torneios
        WHERE id = $1`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return models.Torneio{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Torneio])
}

// Update modifica um torneio existente no banco de dados.
func (r *pgTorneioRepository) Update(ctx context.Context, id int, input models.TorneioInput) (int64, error) {
	query := `
        UPDATE torneios
        SET nome = $1, data_inicio = $2, data_fim = $3, id_esporte = $4, id_cidade = $5, id_estado = $6, id_pais = $7
        WHERE id = $8`
	result, err := r.db.Exec(ctx, query,
		input.Nome, input.DataInicio, input.DataFim, input.EsporteID, input.CidadeID, input.EstadoID, input.PaisID, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// Delete remove um torneio do banco de dados.
func (r *pgTorneioRepository) Delete(ctx context.Context, id int) (int64, error) {
	query := "DELETE FROM torneios WHERE id = $1"
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// InscreverJogador insere uma nova inscrição de jogador/dupla em um torneio.
func (r *pgTorneioRepository) InscreverJogador(ctx context.Context, inscricao models.JogadorTorneio) (models.JogadorTorneio, error) {
	var jogadorInscrito models.JogadorTorneio
	query := `
        INSERT INTO jogadores_torneios (id_torneio, id_jogador, id_categoria, id_dupla, tipo_modalidade)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, id_torneio, id_jogador, id_categoria, id_dupla, tipo_modalidade`
	err := r.db.QueryRow(ctx, query,
		inscricao.TorneioID, inscricao.JogadorID, inscricao.CategoriaID, inscricao.DuplaID, inscricao.TipoModalidade,
	).Scan(
		&jogadorInscrito.ID, &jogadorInscrito.TorneioID, &jogadorInscrito.JogadorID,
		&jogadorInscrito.CategoriaID, &jogadorInscrito.DuplaID, &jogadorInscrito.TipoModalidade,
	)
	return jogadorInscrito, err
}

// ListarInscricoesPorTorneio busca todas as inscrições de um torneio com detalhes dos participantes.
func (r *pgTorneioRepository) ListarInscricoesPorTorneio(ctx context.Context, torneioID int) ([]models.InscricaoDetalhada, error) {
	query := `
		SELECT
			jt.id AS inscricao_id,
			jt.tipo_modalidade,
			-- Detalhes do jogador individual (se modalidade for 'simples')
			j.id AS jogador_id,
			j.nome AS jogador_nome,
			-- Detalhes da dupla (se modalidade for 'duplas')
			d.id AS dupla_id,
			d.nome_dupla,
			-- Detalhes do jogador A da dupla
			ja.id AS jogador_a_id,
			ja.nome AS jogador_a_nome,
			-- Detalhes do jogador B da dupla
			jb.id AS jogador_b_id,
			jb.nome AS jogador_b_nome
		FROM jogadores_torneios jt
		LEFT JOIN jogadores j ON jt.id_jogador = j.id AND jt.tipo_modalidade = 'simples'
		LEFT JOIN duplas d ON jt.id_dupla = d.id AND jt.tipo_modalidade = 'duplas'
		LEFT JOIN jogadores ja ON d.id_jogador_a = ja.id
		LEFT JOIN jogadores jb ON d.id_jogador_b = jb.id
		WHERE jt.id_torneio = $1
		ORDER BY jt.id
	`

	rows, err := r.db.Query(ctx, query, torneioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inscricoes []models.InscricaoDetalhada
	for rows.Next() {
		var inscricao models.InscricaoDetalhada
		var jogadorID, duplaID, jogadorA_ID, jogadorB_ID sql.NullInt64
		var jogadorNome, nomeDupla, jogadorA_Nome, jogadorB_Nome sql.NullString

		err := rows.Scan(
			&inscricao.InscricaoID, &inscricao.TipoModalidade,
			&jogadorID, &jogadorNome,
			&duplaID, &nomeDupla,
			&jogadorA_ID, &jogadorA_Nome,
			&jogadorB_ID, &jogadorB_Nome,
		)
		if err != nil {
			return nil, err
		}

		if inscricao.TipoModalidade == "simples" && jogadorID.Valid {
			inscricao.Jogador = &models.JogadorDetalhes{
				ID:   int(jogadorID.Int64),
				Nome: jogadorNome.String,
			}
		} else if inscricao.TipoModalidade == "duplas" && duplaID.Valid {
			inscricao.Dupla = &models.DuplaDetalhes{
				ID:       int(duplaID.Int64),
				JogadorA: &models.JogadorDetalhes{ID: int(jogadorA_ID.Int64), Nome: jogadorA_Nome.String},
				JogadorB: &models.JogadorDetalhes{ID: int(jogadorB_ID.Int64), Nome: jogadorB_Nome.String},
			}
			if nomeDupla.Valid {
				inscricao.Dupla.NomeDupla = &nomeDupla.String
			}
		}
		inscricoes = append(inscricoes, inscricao)
	}

	return inscricoes, nil
}
