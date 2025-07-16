package models

// EstatisticasJogador representa as estatísticas de um jogador em um grupo.
type EstatisticasJogador struct {
	JogadorID     int    `json:"id_jogador"`
	NomeJogador   string `json:"nome_jogador"`
	SetsGanhos    int    `json:"sets_ganhos"`
	PontosGanhos  int    `json:"pontos_ganhos"`
}

// VencedorGrupo representa um vencedor de um grupo.
type VencedorGrupo struct {
	Posicao      int    `json:"posicao"`
	JogadorID    int    `json:"id_jogador"`
	NomeJogador  string `json:"nome_jogador"`
	Criterio     string `json:"criterio"`
	SetsGanhos   int    `json:"sets_ganhos,omitempty"`
	PontosGanhos int    `json:"pontos_ganhos,omitempty"`
}

// ResultadoVencedores é a estrutura final da resposta JSON.
type ResultadoVencedores struct {
	Vencedores []VencedorGrupo `json:"vencedores"`
}
