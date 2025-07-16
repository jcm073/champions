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

// RankedPlayer armazena o ID e o rating de um jogador.
type RankedPlayer struct {
	ID     int
	Rating int
}

// Validate executa a validação na estrutura CriarGruposInput.
func (c *CriarGruposInput) Validate() error {
	return validation.ValidateStruct(c)
}

// DistributePlayersRanked distribui jogadores com base no rating em grupos de 3 a 5.
func DistributePlayersRanked(players []RankedPlayer) ([][]int, error) {
	totalPlayers := len(players)
	if totalPlayers < 3 {
		return nil, fmt.Errorf("é necessário ter no mínimo 3 jogadores para formar um grupo")
	}

	// Lógica para determinar o número ideal de grupos
	numGroups := 0
	bestNumGroups := 0
	minRemainder := totalPlayers

	for g := 1; g <= totalPlayers/3; g++ {
		if totalPlayers/g >= 3 {
			remainder := totalPlayers % g
			if remainder < minRemainder {
				minRemainder = remainder
				bestNumGroups = g
			}
		}
	}
	numGroups = bestNumGroups
	if numGroups == 0 {
		numGroups = 1 // Fallback para um único grupo se nenhum ideal for encontrado
	}

	groups := make([][]int, numGroups)
	for i := range groups {
		groups[i] = []int{}
	}

	// Distribuição em serpente
	for i, player := range players {
		groupIndex := i % numGroups
		if (i/numGroups)%2 != 0 { // Inverte a direção em linhas ímpares
			groupIndex = numGroups - 1 - groupIndex
		}
		groups[groupIndex] = append(groups[groupIndex], player.ID)
	}

	return groups, nil
}
