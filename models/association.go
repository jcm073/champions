package models

import "competitions/validation"

// EsporteAssociationInput é usado para receber os IDs dos esportes
// ao associar um ou mais esportes a um jogador.
type EsporteAssociationInput struct {
	EsporteIDs []int `json:"esporte_ids" validate:"required,gt=0,dive,gt=0"`
}

// Validate executa as regras de validação para a entrada de EsporteAssociationInput.
func (i *EsporteAssociationInput) Validate() error {
	return validation.ValidateStruct(i)
}
