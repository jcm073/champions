package models

import (
	"competitions/validation"
	"fmt"
)

// Grupo representa um grupo de jogadores em um torneio para uma categoria específica.
type Grupo struct {
	ID          int    `json:"id" db:"id"`
	TorneioID   int    `json:"id_torneio" db:"id_torneio"`
	CategoriaID int    `json:"id_categoria" db:"id_categoria"`
	Nome        string `json:"nome" db:"nome"`
}

// GrupoJogador representa a associação entre um jogador e um grupo.
type GrupoJogador struct {
	GrupoID   int `json:"id_grupo" db:"id_grupo"`
	JogadorID int `json:"id_jogador" db:"id_jogador"`
}

// GrupoComJogadores é uma estrutura para retornar um grupo com a lista de seus jogadores.
type GrupoComJogadores struct {
	Grupo
	Jogadores []Usuario `json:"jogadores"`
}

// CriarGruposInput define os parâmetros para a criação de grupos.
type CriarGruposInput struct {
	CategoriaID int `json:"id_categoria" validate:"required,gt=0"`
}

// Validate executa a validação na estrutura CriarGruposInput.
func (c *CriarGruposInput) Validate() error {
	return validation.ValidateStruct(c)
}

// DistributePlayers distribui os jogadores em grupos de 3 a 5.
func DistributePlayers(playerIDs []int) ([][]int, error) {
	totalPlayers := len(playerIDs)
	if totalPlayers < 3 {
		return nil, fmt.Errorf("é necessário ter no mínimo 3 jogadores para formar um grupo")
	}

	numGroups := 0
	for i := totalPlayers / 5; i >= 0; i-- {
		remainingPlayers := totalPlayers - i*5
		if remainingPlayers%4 == 0 {
			numGroups = i + remainingPlayers/4
			break
		}
		if remainingPlayers%3 == 0 {
			numGroups = i + remainingPlayers/3
			break
		}
	}

	if numGroups == 0 {
		return nil, fmt.Errorf("não é possível dividir %d jogadores em grupos de 3, 4 ou 5", totalPlayers)
	}

	groups := make([][]int, numGroups)
	playerIndex := 0
	for i := 0; i < numGroups; i++ {
		groupSize := 0
		if totalPlayers >= 5 {
			groupSize = 5
		} else if totalPlayers >= 3 {
			groupSize = totalPlayers
		} else {
			// This should not happen due to the initial check, but as a safeguard:
			return nil, fmt.Errorf("erro inesperado na distribuição de jogadores")
		}

		groups[i] = make([]int, groupSize)
		for j := 0; j < groupSize; j++ {
			groups[i][j] = playerIDs[playerIndex]
			playerIndex++
		}
		totalPlayers -= groupSize
	}

	return groups, nil
}