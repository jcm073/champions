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
	ErrEsporteNaoEncontrado     = errors.New("esporte não encontrado")
	ErrEsporteInvalido          = errors.New("um ou mais IDs de esporte fornecidos são inválidos ou não existem")
	ErrUsuarioNaoEncontrado     = errors.New("usuário não encontrado")
	ErrUsuarioInativo           = errors.New("usuário inativo")
	ErrUsuarioJaExiste          = errors.New("já existe um usuário com este e-mail")
	ErrUsuarioInvalido          = errors.New("usuário inválido")
	ErrUsuarioSemPermissao      = errors.New("usuário não tem permissão para esta ação")
)

// UsuarioRepository define a interface para operações de banco de dados de usuários.
// Ela inclui métodos para criar, ler, atualizar e excluir usuários,
// além de associar usuários a esportes e recuperar informações relacionadas.
// A interface é projetada para ser implementada por um repositório específico, como o
// postgresUsuarioRepository, que utiliza o pacote pgx para interagir com um banco de dados PostgreSQL.
// Os métodos retornam erros específicos para facilitar o tratamento de casos comuns,
// como usuários não encontrados, associações já existentes e esportes inválidos.
// A interface também inclui métodos para buscar usuários por ID e e-mail, bem como
// para associar usuários a esportes, permitindo uma flexibilidade maior na manipulação de dados
// relacionados a usuários e suas associações com esportes.
// Além disso, ela permite a atualização de senhas e a recuperação de esportes associados a
// um usuário específico, facilitando a implementação de funcionalidades comuns em aplicações
// que lidam com usuários e suas preferências esportivas.
// A implementação dessa interface deve garantir que as operações sejam realizadas de forma segura e eficiente,
// utilizando transações quando necessário para manter a integridade dos dados.
// A interface também define métodos para buscar esportes associados a um usuário e usuários associados a um esporte,
// permitindo consultas complexas e relacionamentos entre usuários e esportes de forma eficiente.
// A interface é uma abstração que permite a troca de implementações sem afetar o restante
// do código, facilitando testes e manutenção. Ela é essencial para a arquitetura de software,
// pois promove a separação de preocupações e a reutilização de código, permitindo que diferentes
// partes da aplicação interajam com o repositório de usuários de maneira consistente e previs
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
	GetEsportesByUsuario(ctx context.Context, userID int) ([]models.Esporte, error)
	GetUsuariosByEsporte(ctx context.Context, esporteID int) ([]models.Usuario, error)
}

// postgresUsuarioRepository é a implementação concreta do repositório de usuários.
// Ela utiliza o pacote pgx para interagir com um banco de dados PostgreSQL.
// Esta implementação fornece métodos para criar, ler, atualizar e excluir usuários,
// além de associar usuários a esportes e recuperar informações relacionadas.
type postgresUsuarioRepository struct {
	db *pgxpool.Pool
}

// NewUsuarioRepository cria uma nova instância do repositório de usuários.
// Esta função recebe um pool de conexões do pgx e retorna uma instância que implementa a interface UsuarioRepository.
// Ela é responsável por inicializar o repositório com a conexão ao banco de dados,
// permitindo que outras partes da aplicação interajam com os dados de usuários de forma consistente.
// A função é essencial para a configuração do repositório, garantindo que ele esteja pronto para
// realizar operações de banco de dados, como inserções, consultas e atualizações.
// Ao criar uma instância do repositório, você pode usar os métodos definidos na interface
// UsuarioRepository para manipular usuários e suas associações com esportes.
// Isso promove uma separação clara entre a lógica de negócios e a persistência de dados,
// facilitando testes e manutenção do código.
func NewUsuarioRepository(db *pgxpool.Pool) UsuarioRepository {
	return &postgresUsuarioRepository{db: db}
}

// AssociateEsporte associa um jogador (baseado no usuarioID) a um esporte.
// Este método realiza as seguintes etapas:
// 1. Inicia uma transação para garantir a atomicidade das operações.
// 2. Busca o ID do jogador correspondente ao usuarioID dentro da transação.
// 3. Verifica se todos os IDs de esporte fornecidos existem no banco de dados.
// 4. Insere as associações na tabela de junção `jogadores_esportes` usando `UNNEST` para inserção em lote.
// 5. Utiliza `ON CONFLICT (id_jogador, id_esporte) DO NOTHING` para evitar erros se a associação já existir.
// 6. Comita a transação se todas as operações forem bem-sucedidas.
// // Retorna um erro específico se o jogador não for encontrado, se algum esporte for inválido,
// // ou se ocorrer qualquer outro erro durante o processo.
// Esta abordagem garante que as associações sejam feitas de forma segura e eficiente,
// mantendo a integridade dos dados no banco de dados.
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

// GetEsportesByUsuario retorna todos os esportes associados a um usuário (via jogador).
// Este método realiza uma consulta SQL que une as tabelas de esportes, jogadores_esportes e jogadores
// para recuperar os esportes praticados pelo usuário especificado pelo userID.
// Ele utiliza a função pgx.CollectRows para mapear os resultados da consulta para uma
// lista de modelos.Esporte, que é retornada ao chamador.
// Se ocorrer um erro durante a consulta ou o mapeamento, ele é retornado para o chamador,
// permitindo que a lógica de negócios trate o erro adequadamente.
// A consulta SQL utiliza INNER JOINs para garantir que apenas os esportes associados ao jogador
// do usuário sejam retornados, filtrando pelo ID do usuário fornecido.
// Isso garante que a função seja eficiente e retorne apenas os dados necessários,
// evitando a necessidade de carregar dados desnecessários na memória.
// A função é útil para obter rapidamente os esportes que um usuário pratica, permitindo que a
// aplicação apresente essas informações de forma eficiente e organizada.
// Ela é especialmente útil em cenários onde é necessário exibir as preferências esportivas de
// um usuário, como em perfis de usuário ou páginas de configuração de preferências esportivas
func (r *postgresUsuarioRepository) GetEsportesByUsuario(ctx context.Context, userID int) ([]models.Esporte, error) {
	query := `
		SELECT e.id, e.nome, e.descricao
		FROM esportes e
		INNER JOIN jogadores_esportes je ON je.id_esporte = e.id
		INNER JOIN jogadores j ON j.id = je.id_jogador
		WHERE j.id_usuario = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	esportes, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Esporte])
	if err != nil {
		return nil, fmt.Errorf("erro ao mapear esportes do usuário: %w", err)
	}
	return esportes, nil
}

// Create insere um novo usuário no banco de dados.
// Este método realiza uma inserção na tabela 'usuarios' e retorna o ID do usuário recém
// criado, juntamente com a data de criação. A senha deve ser criptografada antes de
// ser passada para este método, garantindo que as senhas sejam armazenadas de forma segura.
// A inserção é feita utilizando uma instrução SQL parametrizada para evitar injeções de
// SQL e garantir a segurança dos dados. O método utiliza o contexto fornecido para
// permitir o cancelamento da operação se necessário, e retorna um erro se a inserção falhar
// por qualquer motivo, como violação de restrições de unicidade (por exemplo,
// se já existir um usuário com o mesmo e-mail).
// A função é útil para criar novos usuários no sistema, permitindo que eles se registrem
// e comecem a interagir com a aplicação. Ela é parte fundamental do fluxo de
// registro de usuários, onde novos usuários fornecem suas informações e são adicionados ao banco de
// dados. A função também pode ser estendida para incluir validações adicionais antes da inserção,
// como verificar se o e-mail já está em uso ou se os dados fornecidos atendem aos
// critérios de formato e consistência esperados.
// A função retorna um erro específico se a inserção falhar, permitindo que a lógica de negócios trate
// o erro de forma adequada, como informando ao usuário que o registro falhou
// ou que o e-mail já está em uso. Isso melhora a experiência do usuário, forne
// cendo feedback claro sobre o que deu errado durante o processo de registro.
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

// FindAll recupera todos os usuários do banco de dados.
// Este método executa uma consulta SQL que seleciona todos os campos da tabela 'usuarios'
// e retorna uma lista de modelos.Usuario. Ele utiliza o contexto fornecido para permitir
// o cancelamento da operação se necessário, e utiliza a função pgx.CollectRows para
// mapear os resultados da consulta para uma lista de usuários.
// A consulta é ordenada pelo nome do usuário para facilitar a leitura e organização dos dados.
// Se ocorrer um erro durante a consulta ou o mapeamento, ele é retornado para o chamador,
// permitindo que a lógica de negócios trate o erro adequadamente.
// A função é útil para obter uma lista completa de usuários registrados no sistema,
// permitindo que a aplicação apresente essas informações de forma eficiente e organizada.
// Ela é especialmente útil em cenários administrativos, onde é necessário visualizar todos os usuários
// registrados, como em painéis de controle ou relatórios de usuários.
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

// FindByID recupera um usuário pelo seu ID.
// Este método executa uma consulta SQL que seleciona todos os campos da tabela 'usuarios'
// onde o ID do usuário corresponde ao fornecido. Ele utiliza o contexto fornecido para
// permitir o cancelamento da operação se necessário.
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

// FindByEmail recupera um usuário pelo seu e-mail.
// Este método executa uma consulta SQL que seleciona todos os campos da tabela 'usuarios'
// onde o e-mail do usuário corresponde ao fornecido. Ele utiliza o contexto fornecido para
// permitir o cancelamento da operação se necessário.
// Se o usuário não for encontrado, ele retorna um erro específico (pgx.ErrNoRows),
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

// FindByIDForAuth recupera um usuário pelo seu ID para autenticação.
// Este método executa uma consulta SQL que seleciona todos os campos relevantes da tabela 'usuarios'
// onde o ID do usuário corresponde ao fornecido. Ele é usado especificamente para autenticação,
// garantindo que os dados necessários para autenticar o usuário estejam disponíveis.
// Ele utiliza o contexto fornecido para permitir o cancelamento da operação se necessário.
// A consulta retorna os campos necessários para autenticação, como ID, tipo, nome, username,
// CPF, data de nascimento, e-mail, senha, telefone, Instagram, data de criação e ativo.
// Se o usuário não for encontrado, ele retorna um erro específico (pgx.ErrNoRows
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

// Update atualiza as informações de um usuário no banco de dados.
// Este método executa uma instrução SQL que atualiza os campos do usuário na tabela 'usuarios'
// com base no ID do usuário fornecido. Ele utiliza o contexto fornecido para permitir o
// cancelamento da operação se necessário.
// A atualização inclui campos como tipo, nome, username, CPF, data de nascimento, e-mail,
// telefone, Instagram e ativo. A senha não é atualizada por este método, pois deve ser
// tratada separadamente através do método UpdatePassword.
// A função retorna o número de linhas afetadas pela atualização, o que pode ser útil para
// verificar se a atualização foi bem-sucedida. Se ocorrer um erro durante a execução da
// instrução SQL, ele é retornado para o chamador, permitindo que a lógica de negócios trate o erro adequadamente.
// É importante garantir que o usuário exista antes de chamar este método,
// para evitar erros de chave estrangeira ou inconsistências no banco de dados.
// A função é útil para atualizar as informações de um usuário existente, permitindo que os usuários
// modifiquem seus dados pessoais, como nome, e-mail, telefone e outras informações relevantes.
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

// UpdatePassword atualiza a senha de um usuário pelo seu ID.
// Retorna o número de linhas afetadas ou um erro, se ocorrer.
// Esta função é útil para permitir que os usuários atualizem suas senhas de forma segura.
// Ela executa uma instrução SQL que atualiza o campo 'password' na tabela 'usuarios'
// onde o ID do usuário corresponde ao fornecido. A senha deve ser previamente criptografada
// antes de ser passada para esta função, garantindo que as senhas sejam armazenadas de forma
func (r *postgresUsuarioRepository) UpdatePassword(ctx context.Context, id uint, newPassword string) (int64, error) {
	res, err := r.db.Exec(ctx, "UPDATE usuarios SET password = $1 WHERE id = $2", newPassword, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

// Delete remove um usuário do banco de dados pelo seu ID.
// Retorna o número de linhas afetadas ou um erro, se ocorrer.
// Esta função é útil para excluir usuários que não são mais necessários ou que solicitaram a exclusão de suas contas.
// Ela executa uma instrução SQL que remove o registro do usuário da tabela 'usuarios'
// onde o ID do usuário corresponde ao fornecido. É importante garantir que o usuário exista
// antes de chamar esta função, para evitar erros de chave estrangeira ou inconsistências no banco
func (r *postgresUsuarioRepository) Delete(ctx context.Context, id int) (int64, error) {
	res, err := r.db.Exec(ctx, "DELETE FROM usuarios WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

// GetUsuariosByEsporte retorna todos os usuários associados a um esporte específico.
// Esta função faz uma junção entre a tabela de usuários e a tabela de associação de usuários
// e esportes, filtrando pelo ID do esporte fornecido.
// Ela retorna uma lista de usuários que praticam o esporte especificado.
func (r *postgresUsuarioRepository) GetUsuariosByEsporte(ctx context.Context, esporteID int) ([]models.Usuario, error) {
	rows, err := r.db.Query(ctx, `
        -- Primeiro, verificar se o esporte existe
        SELECT COUNT(*) FROM esportes WHERE id = $1
    `, esporteID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência do esporte: %w", err)
	}
	var exists int
	rows.Next() // Move to the first (and only) row
	rows.Scan(&exists)
	rows.Close() // Close the first rows object
	if exists == 0 {
		return nil, ErrEsporteNaoEncontrado
	}

	rows, err = r.db.Query(ctx, `
        SELECT u.id, u.tipo, u.nome, u.username, u.cpf, u.data_nascimento, u.email, u.telefone, u.instagram, u.criado_em, u.ativo
        FROM usuarios u
        JOIN usuarios_esportes ue ON u.id = ue.usuario_id
        WHERE ue.esporte_id = $1
    `, esporteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Usuario])
	if err != nil {
		return nil, fmt.Errorf("erro ao mapear usuários por esporte: %w", err)
	}
	return usuarios, nil
}
