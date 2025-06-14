package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Usuario representa um usuário do sistema, com diferentes tipos de acesso e informações pessoais.
// Os tipos de usuário são: jogador, usuario, admin, gestor_clube e gestor_torneio.
// A tabela é criada com o nome "usuarios" e possui os seguintes campos:
// - ID: Identificador único do usuário (chave primária).
// - Tipo: Tipo de usuário, que pode ser 'jogador', 'usuario', 'admin', 'gestor_clube' ou 'gestor_torneio'.
// - Nome: Nome completo do usuário, com tamanho máximo de 100 caracteres.
// - CPF: Cadastro de Pessoa Física, único e obrigatório, com tamanho máximo de 14 caracteres.
// - DataNascimento: Data de nascimento do usuário, obrigatória.
// - Email: Endereço de e-mail do usuário, único e obrigatório, com tamanho máximo de 100 caracteres.
// - Telefone: Número de telefone do usuário, opcional, com tamanho máximo de 20 caracteres.
// - Instagram: Nome de usuário do Instagram, opcional, com tamanho máximo de 50 caracteres.
// - CriadoEm: Data e hora de criação do registro, preenchida automaticamente.
type Usuario struct {
	ID       uint   `json:"id"`
	Tipo     string `json:"tipo" validate:"required,oneof=jogador usuario admin gestor_clube gestor_torneio"`
	Nome     string `json:"nome" validate:"required,max=100"`
	Username string `json:"username" validate:"required,max=50"`
	CPF      string `json:"cpf" validate:"required,max=14"`
	// DataNascimento is stored as a string in "dd/mm/yyyy" format.
	DataNascimento string    `json:"data_nascimento" validate:"required,datetime=year=2000,month=1,day=2"` // Example: "01/01/2000"
	Email          string    `json:"email" validate:"required,email,max=100"`
	Password       string    `json:"password,omitempty" validate:"required,min=8,max=255"`
	Telefone       string    `json:"telefone" validate:"required,min=9,max=20"`
	Instagram      string    `json:"instagram" validate:"max=50"`
	CriadoEm       time.Time `json:"criado_em"`
	Ativo          bool      `json:"ativo"`
}

// Validação usando go-playground/validator

// Validate checks the validation rules for the Usuario struct.
func (u *Usuario) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
