package models

import (
	"competitions/validation"
	"time"
)

// Torneio representa um torneio esportivo, com informações sobre o esporte, nome, datas e localização.
// A estrutura contém campos para o ID do torneio, nome, datas de início e fim,
// IDs de esporte, cidade, estado e país, além de um campo para a data de criação.
// A estrutura é usada para mapear os dados de um torneio em um sistema de competitions.
// Ela pode ser usada para criar, atualizar, buscar e deletar torneios no sistema.
//
//	@Description	Torneio é uma estrutura que representa um torneio esportivo.
//	@ID				Torneio
//	@Name			Torneio
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Torneio
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/torneios [get]
//	@Router			/torneios/{id} [get]
//	@Router			/torneios [post]
//	@Router			/torneios/{id} [put]
//	@Router			/torneios/{id} [delete]
//	@Tags			Torneios
//	@Param			id			path	int				true	"ID do Torneio"
//	@Param			torneio		body	TorneioInput	true	"Dados do Torneio"
//	@Param			nome		query	string			false	"Nome do Torneio"
//	@Param			data_inicio	query	string			false	"Data de Início do Torneio"
//	@Param			data_fim	query	string			false	"Data de Fim do Torneio"
//	@Param			id_esporte	query	int				false	"ID do Esporte"
//	@Param			id_cidade	query	int				false	"ID da Cidade"
//	@Param			id_estado	query	int				false	"ID do Estado"
//	@Param			id_pais		query	int				false	"ID do País"
//	@Param			page		query	int				false	"Número da página para paginação"
//	@Param			limit		query	int				false	"Número de itens por página para paginação"
//	@Param			sort		query	string			false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search		query	string			false	"Termo de busca para filtrar torneios pelo nome"
//	@Security		BearerAuth
type Torneio struct {
	ID         int       `json:"id,omitempty" db:"id"`
	Nome       string    `json:"nome" db:"nome"`
	DataInicio time.Time `json:"data_inicio" db:"data_inicio"`
	DataFim    time.Time `json:"data_fim" db:"data_fim"`
	EsporteID  int       `json:"id_esporte" db:"id_esporte"`
	CidadeID   int       `json:"id_cidade" db:"id_cidade"`
	EstadoID   int       `json:"id_estado" db:"id_estado"`
	PaisID     int       `json:"id_pais" db:"id_pais"`
	CriadoEm   time.Time `json:"criado_em,omitempty" db:"criado_em"`
}

// TorneioInput é usado para receber dados de entrada ao criar ou atualizar um torneio.
// A validação garante que os campos obrigatórios estejam preenchidos e que as datas sejam válidas.
//
//	@Description	TorneioInput é uma estrutura que contém os dados necessários para criar ou atualizar um torneio.
//	@ID				TorneioInput
//	@Name			TorneioInput
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Torneio
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/torneios [post]
//	@Router			/torneios/{id} [put]
//	@Security		BearerAuth
//	@Tags			Torneios
//	@Param			torneio		body	TorneioInput	true	"Dados do Torneio"
//	@Param			nome		query	string			false	"Nome do Torneio"
//	@Param			data_inicio	query	string			false	"Data de Início do Torneio"
//	@Param			data_fim	query	string			false	"Data de Fim do Torneio"
//	@Param			id_esporte	query	int				false	"ID do Esporte"
//	@Param			id_cidade	query	int				false	"ID da Cidade"
//	@Param			id_estado	query	int				false	"ID do Estado"
//	@Param			id_pais		query	int				false	"ID do País"
//	@Param			page		query	int				false	"Número da página para paginação"
//	@Param			limit		query	int				false	"Número de itens por página para paginação"
//	@Param			sort		query	string			false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search		query	string			false	"Termo de busca para filtrar torneios pelo nome"
type TorneioInput struct {
	Nome       string    `json:"nome" validate:"required,max=100"`
	DataInicio time.Time `json:"data_inicio" validate:"required"`
	DataFim    time.Time `json:"data_fim" validate:"required,gtefield=DataInicio"`
	EsporteID  int       `json:"id_esporte" validate:"required,gt=0"`
	CidadeID   int       `json:"id_cidade" validate:"required,gt=0"`
	EstadoID   int       `json:"id_estado" validate:"required,gt=0"`
	PaisID     int       `json:"id_pais" validate:"required,gt=0"`
}

// Validação usando go-playground/validator
func (t *TorneioInput) Validate() error {
	return validation.ValidateStruct(t)
}
