package models

import (
	"competitions/validation"
	"time"
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
	ID             uint      `json:"id,omitempty"`
	Tipo           string    `json:"tipo"`
	Nome           string    `json:"nome"`
	Username       string    `json:"username"`
	CPF            string    `json:"cpf"`
	DataNascimento string    `json:"data_nascimento"`
	Email          string    `json:"email"`
	Password       string    `json:"password,omitempty"`
	Telefone       string    `json:"telefone"`
	Instagram      string    `json:"instagram,omitempty"`
	CriadoEm       time.Time `json:"criado_em,omitempty"`
	Ativo          bool      `json:"ativo"`
}

// UsuarioInput é usado para receber dados de entrada ao criar um usuário.
type UsuarioInput struct {
	Tipo           string `json:"tipo" validate:"required,oneof=jogador usuario admin gestor_clube gestor_torneio"`
	Nome           string `json:"nome" validate:"required,max=100"`
	Username       string `json:"username" validate:"required,max=50"`
	CPF            string `json:"cpf" validate:"required,max=14"`
	DataNascimento string `json:"data_nascimento" validate:"required"` // Formato: "YYYY-MM-DD" para compatibilidade com o tipo DATE do SQL
	Email          string `json:"email" validate:"required,email,max=100"`
	Password       string `json:"password" validate:"required,min=8,max=255"`
	Telefone       string `json:"telefone" validate:"required,min=9,max=20"`
	Instagram      string `json:"instagram" validate:"max=50"`
	Ativo          *bool  `json:"ativo" validate:"required"`
}

// Validate checks the validation rules for the Usuario struct.
func (ui *UsuarioInput) Validate() error {
	return validation.ValidateStruct(ui)
}

// UpdateUsuarioInput é usado para receber dados de entrada ao atualizar um usuário.
// A senha não é incluída aqui; deve ser atualizada através de um endpoint específico.
type UpdateUsuarioInput struct {
	Tipo           string `json:"tipo" validate:"required,oneof=jogador usuario admin gestor_clube gestor_torneio"`
	Nome           string `json:"nome" validate:"required,max=100"`
	Username       string `json:"username" validate:"required,max=50"`
	CPF            string `json:"cpf" validate:"required,max=14"`
	DataNascimento string `json:"data_nascimento" validate:"required"` // Formato: "YYYY-MM-DD"
	Email          string `json:"email" validate:"required,email,max=100"`
	Telefone       string `json:"telefone" validate:"required,min=9,max=20"`
	Instagram      string `json:"instagram" validate:"max=50"`
	Ativo          *bool  `json:"ativo" validate:"required"`
}

func (uui *UpdateUsuarioInput) Validate() error {
	return validation.ValidateStruct(uui)
}

// ChangePasswordInput é usado para receber dados de entrada ao mudar a senha de um usuário.
type ChangePasswordInput struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=255"`
}

// Validate checks the validation rules for the ChangePasswordInput struct.
func (cpi *ChangePasswordInput) Validate() error {
	return validation.ValidateStruct(cpi)
}
