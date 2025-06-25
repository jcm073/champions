package models

import "competitions/validation"

// Esporte representa a entidade 'esportes' no banco de dados.
type Esporte struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

// EsporteInput é usado para receber dados de entrada ao criar ou atualizar um esporte.
// A validação 'oneof' garante que apenas os valores definidos no ENUM do banco de dados sejam aceitos.
type EsporteInput struct {
	Nome string `json:"nome" validate:"required,oneof='Beach Tenis' 'Tenis de Mesa' 'Tenis' 'Pickleball' 'Squash' 'Badminton' 'Padel'"`
}

// Validate executa as regras de validação para a entrada de Esporte.
func (ei *EsporteInput) Validate() error {
	return validation.ValidateStruct(ei)
}
