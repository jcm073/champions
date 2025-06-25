package repository

import (
	"competitions/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UsuarioRepository defines the interface for user data operations.
type UsuarioRepository interface {
	FindAll(ctx context.Context) ([]models.Usuario, error)
	FindByIDForAuth(ctx context.Context, id uint) (*models.Usuario, error)
	FindByID(ctx context.Context, id int) (*models.Usuario, error)
	FindByEmail(ctx context.Context, email string) (*models.Usuario, error)
	Create(ctx context.Context, usuario *models.Usuario) error
	Update(ctx context.Context, id int, input *models.UpdateUsuarioInput) (int64, error)
	Delete(ctx context.Context, id int) (int64, error)
	UpdatePassword(ctx context.Context, userID uint, newHashedPassword string) (int64, error)
}

// pgUsuarioRepository is the PostgreSQL implementation of UsuarioRepository.
type pgUsuarioRepository struct {
	db *pgxpool.Pool
}

// NewUsuarioRepository creates a new instance of UsuarioRepository.
func NewUsuarioRepository(db *pgxpool.Pool) UsuarioRepository {
	return &pgUsuarioRepository{db: db}
}

func (r *pgUsuarioRepository) FindAll(ctx context.Context) ([]models.Usuario, error) {
	query := "SELECT id, tipo, nome, username, cpf, data_nascimento, email, telefone, instagram, criado_em, ativo FROM usuarios ORDER BY nome"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		err := rows.Scan(&u.ID, &u.Tipo, &u.Nome, &u.Username, &u.CPF, &u.DataNascimento, &u.Email, &u.Telefone, &u.Instagram, &u.CriadoEm, &u.Ativo)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}

func (r *pgUsuarioRepository) FindByIDForAuth(ctx context.Context, id uint) (*models.Usuario, error) {
	var u models.Usuario
	// Esta consulta inclui a senha e deve ser usada apenas para fins de autenticação/autorização.
	query := "SELECT id, tipo, nome, username, cpf, data_nascimento, email, password, telefone, instagram, criado_em, ativo FROM usuarios WHERE id=$1"
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Tipo, &u.Nome, &u.Username, &u.CPF, &u.DataNascimento, &u.Email, &u.Password, &u.Telefone, &u.Instagram, &u.CriadoEm, &u.Ativo)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *pgUsuarioRepository) FindByID(ctx context.Context, id int) (*models.Usuario, error) {
	var u models.Usuario
	query := "SELECT id, tipo, nome, username, cpf, data_nascimento, email, telefone, instagram, criado_em, ativo FROM usuarios WHERE id=$1"
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Tipo, &u.Nome, &u.Username, &u.CPF, &u.DataNascimento, &u.Email, &u.Telefone, &u.Instagram, &u.CriadoEm, &u.Ativo)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *pgUsuarioRepository) FindByEmail(ctx context.Context, email string) (*models.Usuario, error) {
	var user models.Usuario
	query := `SELECT id, username, email, password FROM usuarios WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *pgUsuarioRepository) Create(ctx context.Context, usuario *models.Usuario) error {
	query := `
		INSERT INTO usuarios (tipo, nome, username, cpf, data_nascimento, email, password, telefone, instagram, ativo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, criado_em
	`
	return r.db.QueryRow(
		ctx,
		query,
		usuario.Tipo, usuario.Nome, usuario.Username, usuario.CPF, usuario.DataNascimento,
		usuario.Email, usuario.Password, usuario.Telefone, usuario.Instagram, usuario.Ativo,
	).Scan(&usuario.ID, &usuario.CriadoEm)
}

func (r *pgUsuarioRepository) Update(ctx context.Context, id int, input *models.UpdateUsuarioInput) (int64, error) {
	query := `
		UPDATE usuarios SET tipo=$1, nome=$2, username=$3, cpf=$4, data_nascimento=$5, email=$6, telefone=$7, instagram=$8, ativo=$9
		WHERE id=$10
	`
	result, err := r.db.Exec(
		ctx,
		query,
		input.Tipo, input.Nome, input.Username, input.CPF, input.DataNascimento,
		input.Email, input.Telefone, input.Instagram, *input.Ativo, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *pgUsuarioRepository) Delete(ctx context.Context, id int) (int64, error) {
	result, err := r.db.Exec(ctx, "DELETE FROM usuarios WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *pgUsuarioRepository) UpdatePassword(ctx context.Context, userID uint, newHashedPassword string) (int64, error) {
	updateQuery := "UPDATE usuarios SET password = $1 WHERE id = $2"
	result, err := r.db.Exec(ctx, updateQuery, newHashedPassword, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
