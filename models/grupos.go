package models

// Grupo representa um grupo de jogadores/duplas dentro de uma categoria de um torneio.
// Ele pode conter uma lista de JogadoresTorneio, que são os participantes desse grupo.
// A relação entre grupos e jogadores_torneios é N:N, onde um grupo pode ter vários jogadores/duplas e um jogador/dupla pode pertencer a vários grupos.
// A tabela de junção 'grupos_jogadores_torneios' é usada para mapear essa relação.
// // A estrutura 'Grupo' contém o ID do grupo, o ID da categoria a que pertence e o nome do grupo.
// Opcionalmente, pode incluir uma lista de JogadoresTorneio,
// que pode ser preenchida com uma query separada usando a tabela de junção 'grupos_jogadores_torneios'.
//	@Description	Grupo é uma estrutura que representa um grupo de jogadores/duplas dentro de uma categoria de um torneio.
//	@ID				Grupo
//	@Name			Grupo
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Grupo
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/grupos [get]
//	@Router			/grupos/{id} [get]
//	@Router			/grupos [post]
//	@Router			/grupos/{id} [put]
//	@Router			/grupos/{id} [delete]
//	@Tags			Grupos
//	@Param			id		path	int			true	"ID do Grupo"
//	@Param			grupo	body	GrupoInput	true	"Dados do Grupo"
//	@Param			nome	query	string		false	"Nome do Grupo"
//	@Param			page	query	int			false	"Número da página para paginação"
//	@Param			limit	query	int			false	"Número de itens por página para paginação"
//	@Param			sort	query	string		false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search	query	string		false	"Termo de busca para filtrar grupos pelo nome"
//	@Security		BearerAuth
type Grupo struct {
	ID          int    `json:"id"`
	CategoriaID int    `json:"id_categoria"`
	Nome        string `json:"nome"`
	// Opcional: pode ser preenchido com uma query separada usando a tabela de junção.
	JogadoresTorneio []JogadorTorneio `json:"jogadores_torneio,omitempty"`
}

// GrupoJogadorTorneio representa a relação N:N entre 'grupos' e 'jogadores_torneios'.
// Ele é usado para mapear quais jogadores/duplas pertencem a quais grupos.
// A tabela de junção 'grupos_jogadores_torneios' contém os IDs de 'grupos' e 'jogadores_torneios'.
//	@Description	GrupoJogadorTorneio é uma estrutura que representa a relação entre um grupo e um jogador/dupla em um torneio.
//	@ID				GrupoJogadorTorneio
//	@Name			GrupoJogadorTorneio
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	GrupoJogadorTorneio
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/grupos_jogadores_torneios [get]
//	@Router			/grupos_jogadores_torneios/{id} [get]
//	@Router			/grupos_jogadores_torneios [post]
//	@Router			/grupos_jogadores_torneios/{id} [put]
//	@Router			/grupos_jogadores_torneios/{id} [delete]
//	@Tags			GruposJogadoresTorneios
//	@Param			id						path	int							true	"ID da Relação GrupoJogadorTorneio"
//	@Param			grupo_jogador_torneio	body	GrupoJogadorTorneioInput	true	"Dados da Relação GrupoJogadorTorneio"
//	@Param			id_grupo				query	int							false	"ID do Grupo"
//	@Param			id_jogador_torneio		query	int							false	"ID do JogadorTorneio"
//	@Param			page					query	int							false	"Número da página para paginação"
//	@Param			limit					query	int							false	"Número de itens por página para paginação"
//	@Param			sort					query	string						false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search					query	string						false	"Termo de busca para filtrar relações entre grupos e jogadores/duplas"
//	@Security		BearerAuth
type GrupoJogadorTorneio struct {
	GrupoID          int `json:"id_grupo"`
	JogadorTorneioID int `json:"id_jogador_torneio"` // ANTES: id_participante
}
