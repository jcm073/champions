package models

import "competitions/validation"

// Esporte representa um esporte no sistema.
//	@Description	Esporte é uma estrutura que contém informações sobre um esporte específico, incluindo seu ID e nome.
//	@ID				Esporte
//	@Name			Esporte
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Esporte
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/esportes [get]
//	@Router			/esportes/{id} [get]
//	@Router			/esportes [post]
//	@Router			/esportes/{id} [put]
//	@Router			/esportes/{id} [delete]
//	@Tags			Esportes
//	@Param			id		path	int				true	"ID do Esporte"
//	@Param			esporte	body	EsporteInput	true	"Dados do Esporte"
//	@Param			nome	query	string			false	"Nome do Esporte"
//	@Param			page	query	int				false	"Número da página para paginação"
//	@Param			limit	query	int				false	"Número de itens por página para paginação"
//	@Param			sort	query	string			false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search	query	string			false	"Termo de busca para filtrar esportes pelo nome"
//	@Security		BearerAuth
type Esporte struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

// EsporteInput é usado para receber dados de entrada ao criar ou atualizar um esporte.
// A validação 'oneof' garante que apenas os valores definidos no ENUM do banco de dados sejam aceitos.
//	@Description	EsporteInput é uma estrutura que contém os dados necessários para criar ou atualizar um esporte.
//	@ID				EsporteInput
//	@Name			EsporteInput
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Esporte
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/esportes [post]
//	@Router			/esportes/{id} [put]
//	@Tags			Esportes
//	@Param			esporte	body	EsporteInput	true	"Dados do Esporte"
//	@Security		BearerAuth
//	@Param			nome	query	string	false	"Nome do Esporte"
//	@Param			page	query	int		false	"Número da página para paginação"
//	@Param			limit	query	int		false	"Número de itens por página para paginação"
//	@Param			sort	query	string	false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search	query	string	false	"Termo de busca para filtrar esportes pelo nome"
//	@Security		BearerAuth
//	@Param			id		path	int				true	"ID do Esporte"
//	@Param			esporte	body	EsporteInput	true	"Dados do Esporte"
//	@Param			nome	query	string			false	"Nome do Esporte"
//	@Param			page	query	int				false	"Número da página para paginação"
//	@Param			limit	query	int				false	"Número de itens por página para paginação"
//	@Param			sort	query	string			false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search	query	string			false	"Termo de busca para filtrar esportes pelo nome"
//	@Security		BearerAuth
type EsporteInput struct {
	Nome string `json:"nome" validate:"required,oneof='Beach Tenis' 'Tenis de Mesa' 'Tenis' 'Pickleball' 'Squash' 'Badminton' 'Padel'"`
}

// Validate executa as regras de validação para a entrada de Esporte.
func (ei *EsporteInput) Validate() error {
	return validation.ValidateStruct(ei)
}
