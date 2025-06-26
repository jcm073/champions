package repository

import (
	"competitions/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EsporteRepository define a interface para as operações de dados de esportes.
type EsporteRepository interface {
	Create(ctx context.Context, input models.EsporteInput) (*models.Esporte, error)
	FindAll(ctx context.Context) ([]models.Esporte, error)
	FindByID(ctx context.Context, id int) (*models.Esporte, error)
	Update(ctx context.Context, id int, input models.EsporteInput) (int64, error)
	Delete(ctx context.Context, id int) (int64, error)
}

type pgEsporteRepository struct {
	db *pgxpool.Pool
}

// NewEsporteRepository cria uma nova instância de EsporteRepository.
func NewEsporteRepository(db *pgxpool.Pool) EsporteRepository {
	return &pgEsporteRepository{db: db}
}

func (r *pgEsporteRepository) Create(ctx context.Context, input models.EsporteInput) (*models.Esporte, error) {
	var esporte models.Esporte
	query := "INSERT INTO esportes (nome) VALUES ($1) RETURNING id, nome"
	err := r.db.QueryRow(ctx, query, input.Nome).Scan(&esporte.ID, &esporte.Nome)
	if err != nil {
		return nil, err
	}
	return &esporte, nil
}

func (r *pgEsporteRepository) FindAll(ctx context.Context) ([]models.Esporte, error) {
	query := "SELECT id, nome FROM esportes ORDER BY nome"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var esportes []models.Esporte
	for rows.Next() {
		var e models.Esporte
		if err := rows.Scan(&e.ID, &e.Nome); err != nil {
			return nil, err
		}
		esportes = append(esportes, e)
	}
	return esportes, nil
}

func (r *pgEsporteRepository) FindByID(ctx context.Context, id int) (*models.Esporte, error) {
	var esporte models.Esporte
	query := "SELECT id, nome FROM esportes WHERE id = $1"
	err := r.db.QueryRow(ctx, query, id).Scan(&esporte.ID, &esporte.Nome)
	if err != nil {
		return nil, err
	}
	return &esporte, nil
}

func (r *pgEsporteRepository) Update(ctx context.Context, id int, input models.EsporteInput) (int64, error) {
	query := "UPDATE esportes SET nome = $1 WHERE id = $2"
	result, err := r.db.Exec(ctx, query, input.Nome, id)
	return result.RowsAffected(), err
}

func (r *pgEsporteRepository) Delete(ctx context.Context, id int) (int64, error) {
	query := "DELETE FROM esportes WHERE id = $1"
	result, err := r.db.Exec(ctx, query, id)
	return result.RowsAffected(), err
}
