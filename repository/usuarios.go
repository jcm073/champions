package repository

import (
	"competitions/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Erros customizados para o repositório de usuários.
var (
	ErrAssociationAlreadyExists = errors.New("a associação já existe")
	ErrJogadorNaoEncontrado     = errors.New("jogador não encontrado para o usuário especificado")
	ErrEsporteInvalido          = errors.New("um ou mais IDs de esporte fornecidos são inválidos ou não existem")
)

// UsuarioRepository define a interface para operações de banco de dados de usuários.
type UsuarioRepository interface {
	FindAll(ctx context.Context) ([]models.Usuario, error)
	FindByID(ctx context.Context, id int) (*models.Usuario, error)
	FindByIDForAuth(ctx context.Context, id uint) (*models.Usuario, error)
	FindByEmail(ctx context.Context, email string) (*models.Usuario, error)
	Create(ctx context.Context, usuario *models.Usuario) error // This method already accepts *models.Usuario
	Update(ctx context.Context, usuario *models.Usuario) (int64, error)
	UpdatePassword(ctx context.Context, id uint, newPassword string) (int64, error)
	Delete(ctx context.Context, id int) (int64, error)
	AssociateEsporte(ctx context.Context, usuarioID int, esporteIDs []int) error
}

type postgresUsuarioRepository struct {
	db *pgxpool.Pool
}

// NewUsuarioRepository cria uma nova instância do repositório de usuários.
func NewUsuarioRepository(db *pgxpool.Pool) UsuarioRepository {
	return &postgresUsuarioRepository{db: db}
}

// AssociateEsporte associa um jogador (baseado no usuarioID) a um esporte.
func (r *postgresUsuarioRepository) AssociateEsporte(ctx context.Context, usuarioID int, esporteIDs []int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx) // Garante que a transação será revertida em caso de erro

	// 1. Encontrar o ID do jogador correspondente ao ID do usuário dentro da transação.
	var jogadorID int
	err = tx.QueryRow(ctx, "SELECT id FROM jogadores WHERE id_usuario = $1", usuarioID).Scan(&jogadorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrJogadorNaoEncontrado
		}
		return fmt.Errorf("erro ao buscar ID do jogador: %w", err)
	}

	// 2. Verificar se todos os IDs de esporte existem antes de tentar a inserção.
	// Isso evita erros de chave estrangeira e fornece um feedback mais claro ao cliente.
	var count int
	err = tx.QueryRow(ctx, "SELECT count(*) FROM esportes WHERE id = ANY($1)", esporteIDs).Scan(&count)
	if err != nil {
		return fmt.Errorf("erro ao verificar IDs de esporte: %w", err)
	}
	if count != len(esporteIDs) {
		return ErrEsporteInvalido
	}

	// 3. Inserir as associações na tabela de junção usando UNNEST para inserção em lote.
	// ON CONFLICT (id_jogador, id_esporte) DO NOTHING evita erros se a associação já existir.
	query := `
        INSERT INTO jogadores_esportes (id_jogador, id_esporte)
        SELECT $1, unnest($2::int[])
        ON CONFLICT (id_jogador, id_esporte) DO NOTHING;
    `
	_, err = tx.Exec(ctx, query, jogadorID, esporteIDs)
	if err != nil {
		return fmt.Errorf("erro ao associar esportes: %w", err)
	}

	return tx.Commit(ctx) // Confirma a transação
}

// Implementações dos outros métodos da interface (Create, FindAll, etc.)
// ... (O restante do seu código do repositório de usuário iria aqui)
func (r *postgresUsuarioRepository) Create(ctx context.Context, usuario *models.Usuario) error {
	// A trigger no banco de dados irá popular a tabela 'jogadores' se o tipo for 'jogador'.
	query := `
        INSERT INTO usuarios (tipo, nome, username, cpf, data_nascimento, email, password, telefone, instagram, ativo)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, criado_em`
	return r.db.QueryRow(ctx, query,
		usuario.Tipo, usuario.Nome, usuario.Username, usuario.CPF, usuario.DataNascimento,
		usuario.Email, usuario.Password, usuario.Telefone, usuario.Instagram, usuario.Ativo,
	).Scan(&usuario.ID, &usuario.CriadoEm)
}

func (r *postgresUsuarioRepository) FindAll(ctx context.Context) ([]models.Usuario, error) {
	rows, err := r.db.Query(ctx, "SELECT id, tipo, nome, username, cpf, data_nascimento, email, telefone, instagram, criado_em, ativo FROM usuarios ORDER BY nome")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Usuario])
	if err != nil {
		return nil, fmt.Errorf("erro ao mapear usuários: %w", err)
	}
	return usuarios, nil
}

func (r *postgresUsuarioRepository) FindByID(ctx context.Context, id int) (*models.Usuario, error) {
	row, err := r.db.Query(ctx, "SELECT id, tipo, nome, username, cpf, data_nascimento, email, telefone, instagram, criado_em, ativo FROM usuarios WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	usuario, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Usuario])
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *postgresUsuarioRepository) FindByEmail(ctx context.Context, email string) (*models.Usuario, error) {
	query := "SELECT id, tipo, nome, username, cpf, data_nascimento, email, password, telefone, instagram, criado_em, ativo FROM usuarios WHERE email = $1"
	rows, err := r.db.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuario, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Usuario])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Retorna um erro específico para "não encontrado" que pode ser tratado no handler.
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("erro ao mapear usuário por e-mail: %w", err)
	}
	return &usuario, nil
}

func (r *postgresUsuarioRepository) FindByIDForAuth(ctx context.Context, id uint) (*models.Usuario, error) {
	row, err := r.db.Query(ctx, "SELECT id, tipo, nome, username, cpf, data_nascimento, email, password, telefone, instagram, criado_em, ativo FROM usuarios WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	usuario, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Usuario])
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *postgresUsuarioRepository) Update(ctx context.Context, usuario *models.Usuario) (int64, error) {
	// A trigger no banco de dados irá sincronizar com a tabela 'jogadores'.
	res, err := r.db.Exec(ctx, `
        UPDATE usuarios SET tipo=$1, nome=$2, username=$3, cpf=$4, data_nascimento=$5, email=$6, telefone=$7, instagram=$8, ativo=$9
        WHERE id=$10`,
		usuario.Tipo, usuario.Nome, usuario.Username, usuario.CPF, usuario.DataNascimento, usuario.Email, usuario.Telefone, usuario.Instagram, usuario.Ativo, usuario.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (r *postgresUsuarioRepository) UpdatePassword(ctx context.Context, id uint, newPassword string) (int64, error) {
	res, err := r.db.Exec(ctx, "UPDATE usuarios SET password = $1 WHERE id = $2", newPassword, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (r *postgresUsuarioRepository) Delete(ctx context.Context, id int) (int64, error) {
	res, err := r.db.Exec(ctx, "DELETE FROM usuarios WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
