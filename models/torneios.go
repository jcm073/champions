package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Torneio representa um torneio esportivo, com informações sobre o esporte, nome, descrição e quantidade de quadras.
// A tabela é criada com o nome "torneios" e possui os seguintes campos:
//   - ID: Identificador único do torneio (chave primária).
//   - EsporteID: Identificador do esporte relacionado ao torneio, obrigatório.
//   - Nome: Nome do torneio, obrigatório, com tamanho máximo de 100 caracteres.
//   - Descricao: Descrição do torneio, opcional, com tamanho máximo de 255 caracteres.
//   - QuantidadeQuadras: Número de quadras disponíveis para o torneio,
//     obrigatório, com valor mínimo de 1.
//   - CriadoEm: Data e hora de criação do registro, preenchida automaticamente.
type Torneio struct {
	ID                int       `json:"id,omitempty"`
	EsporteID         int       `json:"id_esporte" validate:"gt=0"`
	Nome              string    `json:"nome" validate:"required,max=100"`
	Descricao         string    `json:"descricao" validate:"max=255"`
	QuantidadeQuadras int       `json:"quantidade_quadras" validate:"min=1"`
	CriadoEm          time.Time `json:"criado_em,omitempty"`
}

// TorneioInput é usado para receber dados de entrada ao criar ou atualizar um torneio.
// Ele é usado para validar os dados antes de criar ou atualizar um registro na tabela "torneios".
type TorneioInput struct {
	EsporteID         int    `json:"id_esporte" validate:"gt=0"`
	Nome              string `json:"nome" validate:"required,max=100"`
	Descricao         string `json:"descricao" validate:"max=255"`
	QuantidadeQuadras int    `json:"quantidade_quadras" validate:"min=1"`
}

// Validação usando go-playground/validator
func (t *TorneioInput) Validate() error {
	return validate.Struct(t)
}
