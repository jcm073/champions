package models

import "time"

// Jogo representa uma partida dentro de um torneio, correspondendo à tabela 'jogos'.
// Os campos foram atualizados para refletir a mudança de 'participantes' para 'jogadores_torneios'.
// A estrutura contém informações sobre o jogo, como IDs dos torneios, grupos, rodadas,
// jogadores/duplas participantes, vencedor, perdedor, tipo de modalidade, data/hora,
// localização, situação do jogo e se é a final do campeonato.
// A estrutura também inclui campos opcionais para os IDs dos jogadores/duplas vencedoras e perdedoras,
// que podem ser nulos se o jogo ainda não tiver sido decidido.
// A estrutura 'Jogo' é usada para mapear os dados de uma partida em um torneio.
// Ela pode ser usada para criar, atualizar, buscar e deletar jogos no sistema.
// A estrutura inclui campos para os IDs dos torneios, grupos e rodadas, bem como
// os IDs dos jogadores/duplas participantes, vencedor, perdedor, tipo de modalidade,
// data/hora do jogo, localização, situação do jogo e se é a final do campeonato.
// que podem ser nulos se o jogo ainda não tiver sido decidido.
//	@Description	Jogo é uma estrutura que representa uma partida dentro de um torneio.
//	@ID				Jogo
//	@Name			Jogo
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Jogo
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/jogos [get]
//	@Router			/jogos/{id} [get]
//	@Router			/jogos [post]
//	@Router			/jogos/{id} [put]
//	@Router			/jogos/{id} [delete]
//	@Tags			Jogos
//	@Param			id					path	int			true	"ID do Jogo"
//	@Param			jogo				body	JogoInput	true	"Dados do Jogo"
//	@Param			id_torneio			query	int			false	"ID do Torneio"
//	@Param			id_grupo			query	int			false	"ID do Grupo"
//	@Param			id_rodada			query	int			false	"ID da Rodada"
//	@Param			id_jogador_torneio1	query	int			false	"ID do JogadorTorneio 1"
//	@Param			id_jogador_torneio2	query	int			false	"ID do JogadorTorneio 2"
//	@Param			id_dupla1			query	int			false	"ID da Dupla 1"
//	@Param			id_dupla2			query	int			false	"ID da Dupla 2"
//	@Param			id_jogador_vencedor	query	int			false	"ID do Jogador Vencedor"
//	@Param			id_jogador_perdedor	query	int			false	"ID do Jogador Perdedor"
//	@Param			id_dupla_vencedora	query	int			false	"ID da Dupla Vencedora"
//	@Param			id_dupla_perdedora	query	int			false	"ID da Dupla Perdedor"
//	@Param			tipo_modalidade		query	string		false	"Tipo de Modalidade"
//	@Param			data_hora			query	string		false	"Data e Hora do Jogo"
//	@Param			localizacao			query	string		false	"Localização do Jogo"
//	@Param			situacao			query	string		false	"Situação do Jogo"
//	@Param			eh_final_campeonato	query	bool		false	"É Final do Campeonato"
//	@Security		BearerAuth
type Jogo struct {
	ID                int       `json:"id"`
	TorneioID         int       `json:"id_torneio"`
	GrupoID           int       `json:"id_grupo"`
	RodadaID          int       `json:"id_rodada"`
	JogadorTorneio1ID *int      `json:"id_jogador_torneio1,omitempty"` // ANTES: id_participante1
	JogadorTorneio2ID *int      `json:"id_jogador_torneio2,omitempty"` // ANTES: id_participante2
	Dupla1ID          *int      `json:"id_dupla1,omitempty"`
	Dupla2ID          *int      `json:"id_dupla2,omitempty"`
	JogadorVencedorID *int      `json:"id_jogador_vencedor,omitempty"`
	JogadorPerdedorID *int      `json:"id_jogador_perdedor,omitempty"`
	DuplaVencedoraID  *int      `json:"id_dupla_vencedora,omitempty"`
	DuplaPerdedoraID  *int      `json:"id_dupla_perdedora,omitempty"`
	TipoModalidade    string    `json:"tipo_modalidade"`
	DataHora          time.Time `json:"data_hora"`
	Localizacao       *string   `json:"localizacao,omitempty"`
	Situacao          string    `json:"situacao"`
	EhFinalCampeonato bool      `json:"eh_final_campeonato"`
}
