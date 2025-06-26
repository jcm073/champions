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
	DataNascimento time.Time `json:"data_nascimento"`
	Email          string    `json:"email"`
	Password       string    `json:"password,omitempty"`
	Telefone       string    `json:"telefone"`
	Instagram      string    `json:"instagram,omitempty"`
	CriadoEm       time.Time `json:"criado_em,omitempty"`
	Ativo          bool      `json:"ativo"`
}

// UsuarioInput é usado para receber dados de entrada ao criar um usuário.
// Ele inclui validações para garantir que os campos obrigatórios estejam preenchidos e que os formatos sejam válidos.
// A validação 'user_type' garante que o tipo de usuário seja um dos valores permitidos: 'jogador', 'usuario', 'admin', 'gestor_clube' ou 'gestor_torneio'.
// A validação 'max' garante que os campos de texto não excedam os limites de tamanho especificados.
// A validação 'required' garante que os campos obrigatórios estejam preenchidos.
// A validação 'email' garante que o campo de e-mail esteja no formato correto.
// A validação 'min' garante que a senha tenha pelo menos 8 caracteres.
// O campo 'ativo' é um ponteiro para bool, permitindo que seja nulo, o que é útil para indicar se o usuário está ativo ou inativo.
// O campo 'data_nascimento' deve ser fornecido no formato "YYYY-MM-DD" para compatibilidade com o tipo DATE do SQL.
// O campo 'telefone' deve ter pelo menos 9 caracteres e no máximo 20,
// enquanto o campo 'instagram' é opcional e pode ter até 50 caracteres.
// O campo 'username' é obrigatório e deve ter no máximo 50 caracteres.
// O campo 'password' é obrigatório e deve ter entre 8 e 255 caracteres.
// O campo 'cpf' é obrigatório e deve ter no máximo 14 caracteres, seguindo o formato brasileiro de CPF (apenas números, com ou sem máscara).
// A estrutura 'UsuarioInput` é usada para mapear os dados de entrada ao criar um novo usuário no sistema.
// Ela pode ser usada em endpoints de criação de usuários, onde os dados são recebidos no corpo da requisição
// e validados antes de serem persistidos no banco de dados.
// A estrutura inclui campos para o tipo de usuário, nome, username, CPF, data de nascimento, e-mail, senha, telefone, Instagram,
// ativo e a data de criação do registro.
// A validação é feita usando a biblioteca 'validation' para garantir que os dados atendam aos requisitos especificados.
// A estrutura é usada para mapear os dados de entrada ao criar um novo usuário no sistema.
//	@Description	UsuarioInput é uma estrutura que contém os dados necessários para criar um novo usuário.
//	@ID				UsuarioInput
//	@Name			UsuarioInput
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Usuario
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/usuarios [post]
//	@Tags			Usuarios
//	@Param			usuario	body	UsuarioInput	true	"Dados do Usuário"
//	@Security		BearerAuth
//	@Param			tipo			query	string	false	"Tipo de Usuário (jogador, usuario, admin, gestor_clube, gestor_torneio)"
//	@Param			nome			query	string	false	"Nome do Usuário"
//	@Param			username		query	string	false	"Username do Usuário"
//	@Param			cpf				query	string	false	"CPF do Usuário"
//	@Param			data_nascimento	query	string	false	"Data de Nascimento do Usuário (YYYY-MM-DD)"
//	@Param			email			query	string	false	"E-mail do Usuário"
//	@Param			telefone		query	string	false	"Telefone do Usuário"
//	@Param			instagram		query	string	false	"Instagram do Usuário"
//	@Param			ativo			query	bool	false	"Ativo do Usuário"
//	@Param			page			query	int		false	"Número da página para paginação"
//	@Param			limit			query	int		false	"Número de itens por página para paginação"
//	@Param			sort			query	string	false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search			query	string	false	"Termo de busca para filtrar usuários pelo nome ou username"
//	@Security		BearerAuth
type UsuarioInput struct {
	Tipo           string    `json:"tipo" validate:"required,user_type"`
	Nome           string    `json:"nome" validate:"required,max=100"`
	Username       string    `json:"username" validate:"required,max=50"`
	CPF            string    `json:"cpf" validate:"required,max=14"`
	DataNascimento time.Time `json:"data_nascimento" validate:"required"` // Formato: "YYYY-MM-DD" para compatibilidade com o tipo DATE do SQL
	Email          string    `json:"email" validate:"required,email,max=100"`
	Password       string    `json:"password" validate:"required,min=8,max=255"`
	Telefone       string    `json:"telefone" validate:"required,min=9,max=20"`
	Instagram      string    `json:"instagram" validate:"max=50"`
	Ativo          *bool     `json:"ativo" validate:"required"`
}

// Validate checks the validation rules for the Usuario struct.
func (ui *UsuarioInput) Validate() error {
	return validation.ValidateStruct(ui)
}

// UpdateUsuarioInput é usado para receber dados de entrada ao atualizar um usuário.
// A senha não é incluída aqui; deve ser atualizada através de um endpoint específico.
// A validação é semelhante à de UsuarioInput, mas não inclui o campo 'password'.
// A estrutura 'UpdateUsuarioInput' é usada para mapear os dados de entrada ao atualizar um usuário existente no sistema.
// Ela pode ser usada em endpoints de atualização de usuários, onde os dados são recebidos no corpo da requisição
// e validados antes de serem persistidos no banco de dados.
// A estrutura inclui campos para o tipo de usuário, nome, username, CPF, data de nascimento, e-mail, telefone, Instagram e ativo.
// A validação é feita usando a biblioteca 'validation' para garantir que os dados atendam aos requisitos especificados.
// A estrutura é usada para mapear os dados de entrada ao atualizar um usuário existente no sistema.
//	@Description	UpdateUsuarioInput é uma estrutura que contém os dados necessários para atualizar um usuário existente.
//	@ID				UpdateUsuarioInput
//	@Name			UpdateUsuarioInput
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Usuario
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/usuarios/{id} [put]
//	@Tags			Usuarios
//	@Param			id		path	int					true	"ID do Usuário"
//	@Param			usuario	body	UpdateUsuarioInput	true	"Dados do Usuário"
//	@Security		BearerAuth
//	@Param			tipo			query	string	false	"Tipo de Usuário (jogador, usuario, admin, gestor_clube, gestor_torneio)"
//	@Param			nome			query	string	false	"Nome do Usuário"
//	@Param			username		query	string	false	"Username do Usuário"
//	@Param			cpf				query	string	false	"CPF do Usuário"
//	@Param			data_nascimento	query	string	false	"Data de Nascimento do Usuário (YYYY-MM-DD)"
//	@Param			email			query	string	false	"E-mail do Usuário"
//	@Param			telefone		query	string	false	"Telefone do Usuário"
//	@Param			instagram		query	string	false	"Instagram do Usuário"
//	@Param			ativo			query	bool	false	"Ativo do Usuário"
//	@Param			page			query	int		false	"Número da página para paginação"
//	@Param			limit			query	int		false	"Número de itens por página para paginação"
//	@Param			sort			query	string	false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search			query	string	false	"Termo de busca para filtrar usuários pelo nome ou username"
//	@Security		BearerAuth
type UpdateUsuarioInput struct {
	Tipo           string    `json:"tipo" validate:"required,oneof=jogador usuario admin gestor_clube gestor_torneio"`
	Nome           string    `json:"nome" validate:"required,max=100"`
	Username       string    `json:"username" validate:"required,max=50"`
	CPF            string    `json:"cpf" validate:"required,max=14"`
	DataNascimento time.Time `json:"data_nascimento" validate:"required"` // Formato: "YYYY-MM-DD"
	Email          string    `json:"email" validate:"required,email,max=100"`
	Telefone       string    `json:"telefone" validate:"required,min=9,max=20"`
	Instagram      string    `json:"instagram" validate:"max=50"`
	Ativo          *bool     `json:"ativo" validate:"required"`
}

func (uui *UpdateUsuarioInput) Validate() error {
	return validation.ValidateStruct(uui)
}

// ChangePasswordInput é usado para receber dados de entrada ao mudar a senha de um usuário.
// Ele inclui validações para garantir que os campos obrigatórios estejam preenchidos e que a nova senha atenda aos requisitos de segurança.
// A validação 'required' garante que os campos obrigatórios estejam preenchidos.
// A validação 'min' garante que a nova senha tenha pelo menos 8 caracteres.
// A validação 'max' garante que a nova senha não exceda 255 caracteres.
// A estrutura 'ChangePasswordInput' é usada para mapear os dados de entrada ao mudar a senha de um usuário no sistema.
// Ela pode ser usada em endpoints de mudança de senha, onde os dados são recebidos no corpo da requisição
// e validados antes de serem persistidos no banco de dados.
// A estrutura inclui campos para a senha antiga e a nova senha.
// A validação é feita usando a biblioteca 'validation' para garantir que os dados atendam
// aos requisitos especificados.
// A estrutura é usada para mapear os dados de entrada ao mudar a senha de um usuário
//	@Description	ChangePasswordInput é uma estrutura que contém os dados necessários para mudar a senha de um usuário.
//	@ID				ChangePasswordInput
//	@Name			ChangePasswordInput
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Usuario
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/usuarios/{id}/change-password [put]
//	@Tags			Usuarios
//	@Param			id				path	int					true	"ID do Usuário"
//	@Param			change_password	body	ChangePasswordInput	true	"Dados para Mudar a Senha"
//	@Security		BearerAuth
type ChangePasswordInput struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=255"`
}

// Validate checks the validation rules for the ChangePasswordInput struct.
func (cpi *ChangePasswordInput) Validate() error {
	return validation.ValidateStruct(cpi)
}
