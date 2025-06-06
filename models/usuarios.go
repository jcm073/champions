package models

import "time"

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
	ID             uint      `gorm:"primaryKey" json:"id"`
	Tipo           string    `gorm:"type:enum('jogador','usuario','admin','gestor_clube','gestor_torneio');not null" json:"tipo"`
	Nome           string    `gorm:"type:varchar(100);not null" json:"nome"`
	Username       string    `gorm:"type:varchar(50);unique;not null" json:"username"`
	CPF            string    `gorm:"type:varchar(14);unique;not null" json:"cpf"`
	DataNascimento time.Time `gorm:"type:datetime;not null" json:"data_nascimento"`
	Email          string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password       string    `gorm:"type:varchar(255);not null" json:"password"`
	Telefone       string    `gorm:"type:varchar(20)" json:"telefone"`
	Instagram      string    `gorm:"type:varchar(50)" json:"instagram"`
	CriadoEm       time.Time `gorm:"autoCreateTime" json:"criadoem"`
	Ativo          bool      `gorm:"type:tinyint(1);default:1" json:"ativo"`
}
