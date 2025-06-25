package repository

import (
	"competitions/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TorneioRepository defines the interface for tournament data operations.
type TorneioRepository interface {
	Create(ctx context.Context, input models.TorneioInput) (*models.Torneio, error)
	FindAll(ctx context.Context) ([]models.Torneio, error)
	FindByID(ctx context.Context, id int) (*models.Torneio, error)
	Update(ctx context.Context, id int, input models.TorneioInput) (int64, error)
	Delete(ctx context.Context, id int) (int64, error)
}

// pgTorneioRepository is the PostgreSQL implementation of TorneioRepository.
type pgTorneioRepository struct {
	db *pgxpool.Pool
}

// NewTorneioRepository creates a new instance of TorneioRepository.
func NewTorneioRepository(db *pgxpool.Pool) TorneioRepository {
	return &pgTorneioRepository{db: db}
}

func (r *pgTorneioRepository) Create(ctx context.Context, input models.TorneioInput) (*models.Torneio, error) {
	var torneio models.Torneio
	query := `
        INSERT INTO torneios (id_esporte, nome, descricao, quantidade_quadras)
        VALUES ($1, $2, $3, $4)
        RETURNING id, id_esporte, nome, descricao, quantidade_quadras, criado_em
    `
	err := r.db.QueryRow(
		ctx, query,
		input.EsporteID, input.Nome, input.Descricao, input.QuantidadeQuadras,
	).Scan(&torneio.ID, &torneio.EsporteID, &torneio.Nome, &torneio.Descricao, &torneio.QuantidadeQuadras, &torneio.CriadoEm)
	if err != nil {
		return nil, err
	}
	return &torneio, nil
}

func (r *pgTorneioRepository) FindAll(ctx context.Context) ([]models.Torneio, error) {
	query := "SELECT id, id_esporte, nome, descricao, quantidade_quadras, criado_em FROM torneios ORDER BY nome"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var torneios []models.Torneio
	for rows.Next() {
		var t models.Torneio
		if err := rows.Scan(&t.ID, &t.EsporteID, &t.Nome, &t.Descricao, &t.QuantidadeQuadras, &t.CriadoEm); err != nil {
			return nil, err
		}
		torneios = append(torneios, t)
	}
	return torneios, nil
}

func (r *pgTorneioRepository) FindByID(ctx context.Context, id int) (*models.Torneio, error) {
	var torneio models.Torneio
	query := "SELECT id, id_esporte, nome, descricao, quantidade_quadras, criado_em FROM torneios WHERE id = $1"
	err := r.db.QueryRow(ctx, query, id).Scan(&torneio.ID, &torneio.EsporteID, &torneio.Nome, &torneio.Descricao, &torneio.QuantidadeQuadras, &torneio.CriadoEm)
	if err != nil {
		return nil, err
	}
	return &torneio, nil
}

func (r *pgTorneioRepository) Update(ctx context.Context, id int, input models.TorneioInput) (int64, error) {
	query := "UPDATE torneios SET id_esporte = $1, nome = $2, descricao = $3, quantidade_quadras = $4 WHERE id = $5"
	result, err := r.db.Exec(ctx, query, input.EsporteID, input.Nome, input.Descricao, input.QuantidadeQuadras, id)
	return result.RowsAffected(), err
}

func (r *pgTorneioRepository) Delete(ctx context.Context, id int) (int64, error) {
	query := "DELETE FROM torneios WHERE id = $1"
	result, err := r.db.Exec(ctx, query, id)
	return result.RowsAffected(), err
}
