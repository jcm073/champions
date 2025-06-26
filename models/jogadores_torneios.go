package models

import "competitions/validation"

// JogadorTorneio representa a inscrição de um jogador ou dupla em um torneio,
// correspondendo à tabela 'jogadores_torneios'.
// Ele contém o ID do torneio, o ID do jogador (ou nulo se for uma dupla),
// o ID da categoria, o ID da dupla (ou nulo se for um jogador individual) e o tipo de modalidade
// (simples ou duplas).
// A estrutura é usada para validar os dados de entrada ao inscrever um jogador ou dupla em
// um torneio, garantindo que os campos obrigatórios sejam preenchidos corretamente.
// A validação é feita usando a biblioteca 'validation' para garantir que os dados estejam corretos antes de serem persistidos no banco de dados.
// A estrutura também possui métodos para converter entre a estrutura de entrada e o modelo de banco de dados,
// facilitando a manipulação dos dados na aplicação
// e garantindo que as regras de negócio sejam respeitadas.
// A estrutura JogadorTorneio é essencial para o gerenciamento de inscrições em torneios,
// permitindo que a aplicação trate tanto jogadores individuais quanto duplas de forma consistente.
// A estrutura também é usada para gerar respostas detalhadas de inscrição, incluindo informações sobre o jogador
// ou dupla inscrita, facilitando a visualização e o gerenciamento das inscrições pelos administradores
// e organizadores de torneios.
//	@Description	JogadorTorneio é uma estrutura que representa a inscrição de um jogador ou dupla em um torneio.
//	@ID				JogadorTorneio
//	@Name			JogadorTorneio
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	JogadorTorneio
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/jogadores_torneios [get]
//	@Router			/jogadores_torneios/{id} [get]
//	@Router			/jogadores_torneios [post]
//	@Router			/jogadores_torneios/{id} [put]
//	@Router			/jogadores_torneios/{id} [delete]
//	@Tags			JogadoresTorneios
//	@Param			id				path	int					true	"ID do JogadorTorneio"
//	@Param			jogador_torneio	body	JogadorTorneioInput	true	"Dados do JogadorTorneio"
//	@Param			torneio_id		query	int					false	"ID do Torneio"
//	@Param			id_jogador		query	int					false	"ID do Jogador"
//	@Param			id_dupla		query	int					false	"ID da Dupla"
//	@Param			id_categoria	query	int					false	"ID da Categoria"
//	@Param			tipo_modalidade	query	string				false	"Tipo de Modalidade (simples ou duplas)"
//	@Param			page			query	int					false	"Número da página para paginação"
//	@Param			limit			query	int					false	"Número de itens por página para paginação"
//	@Param			sort			query	string				false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search			query	string				false	"Termo de busca para filtrar inscrições por jogador, dupla ou categoria"
//	@Security		BearerAuth
//	@Security		JWTAuth
type JogadorTorneio struct {
	ID             int    `json:"id"`
	TorneioID      int    `json:"id_torneio"`
	JogadorID      *int   `json:"id_jogador,omitempty"` // Ponteiro para permitir nulo
	CategoriaID    int    `json:"id_categoria"`
	DuplaID        *int   `json:"id_dupla,omitempty"` // Ponteiro para permitir nulo
	TipoModalidade string `json:"tipo_modalidade"`
}

// JogadorTorneioInput é a estrutura para validar os dados de entrada ao inscrever
// um jogador ou dupla em um torneio.
// Ela contém os campos necessários para a inscrição, incluindo o ID do torneio,
// o ID do jogador (ou nulo se for uma dupla), o ID da categoria,
// o ID da dupla (ou nulo se for um jogador individual) e o tipo de modalidade
// (simples ou duplas).
// A validação é feita usando a biblioteca 'validation' para garantir que os dados estejam corretos antes de serem persistidos no banco de dados.
// A estrutura JogadorTorneioInput é usada para receber os dados de entrada da API
// ao inscrever um jogador ou dupla em um torneio.
// Ela garante que os campos obrigatórios sejam preenchidos corretamente e que o tipo de modalidade
// seja válido (simples ou duplas).
// A estrutura também possui métodos para converter entre a estrutura de entrada e o modelo de banco de dados,
// facilitando a manipulação dos dados na aplicação e garantindo que as regras de negócio sejam respeitadas.
// A estrutura JogadorTorneioInput é essencial para o gerenciamento de inscrições em torneios,
// permitindo que a aplicação trate tanto jogadores individuais quanto duplas de forma consistente.
// Ela é usada para validar os dados de entrada antes de criar ou atualizar uma inscrição,
// garantindo que os dados estejam corretos e completos.
// A estrutura também é usada para gerar respostas detalhadas de inscrição, incluindo informações sobre o jogador
// ou dupla inscrita, facilitando a visualização e o gerenciamento das inscrições pelos administradores
// e organizadores de torneios.
//	@Description	JogadorTorneioInput é uma estrutura que contém os dados necessários para inscrever um jogador ou dupla em um torneio.
//	@ID				JogadorTorneioInput
//	@Name			JogadorTorneioInput
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	JogadorTorneioInput
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/jogadores_torneios [post]
//	@Router			/jogadores_torneios/{id} [put]
//	@Tags			JogadoresTorneios
//	@Param			jogador_torneio	body	JogadorTorneioInput	true	"Dados do JogadorTorneioInput"
//	@Param			id				path	int					true	"ID do JogadorTorneioInput"
//	@Param			torneio_id		query	int					false	"ID do Torneio"
//	@Param			id_jogador		query	int					false	"ID do Jogador"
//	@Param			id_dupla		query	int					false	"ID da Dupla"
//	@Param			id_categoria	query	int					false	"ID da Categoria"
//	@Param			tipo_modalidade	query	string				false	"Tipo de Modalidade (simples ou duplas)"
//	@Security		BearerAuth
//	@Security		JWTAuth
type JogadorTorneioInput struct {
	TorneioID      int    `json:"torneio_id"`
	JogadorID      *int   `json:"id_jogador" validate:"required_if=TipoModalidade simples,excluded_with=DuplaID"`
	CategoriaID    int    `json:"id_categoria" validate:"required,gt=0"`
	DuplaID        *int   `json:"id_dupla" validate:"required_if=TipoModalidade duplas,excluded_with=JogadorID"`
	TipoModalidade string `json:"tipo_modalidade" validate:"required,oneof=simples duplas"`
}

// Validate executa as regras de validação na estrutura de entrada.
func (jti *JogadorTorneioInput) Validate() error {
	return validation.ValidateStruct(jti)
}

// IsValidTipoModalidade verifica se o tipo de modalidade é válido.
func IsValidTipoModalidade(tipo string) bool {
	validTypes := []string{"simples", "duplas"}
	for _, validType := range validTypes {
		if tipo == validType {
			return true
		}
	}
	return false
}

// ToModel converte JogadorTorneioInput para JogadorTorneio.
func (jti *JogadorTorneioInput) ToModel() JogadorTorneio {
	return JogadorTorneio{
		TorneioID:      jti.TorneioID,
		JogadorID:      jti.JogadorID,
		CategoriaID:    jti.CategoriaID,
		DuplaID:        jti.DuplaID,
		TipoModalidade: jti.TipoModalidade,
	}
}

// ToInput converte JogadorTorneio para JogadorTorneioInput.
func (jt *JogadorTorneio) ToInput() JogadorTorneioInput {
	return JogadorTorneioInput{
		TorneioID:      jt.TorneioID,
		JogadorID:      jt.JogadorID,
		CategoriaID:    jt.CategoriaID,
		DuplaID:        jt.DuplaID,
		TipoModalidade: jt.TipoModalidade,
	}
}

// =============================================================================
// Structs para Respostas Detalhadas de Inscrição
// =============================================================================

// JogadorDetalhes contém informações simplificadas de um jogador para a resposta da API.
// Ele inclui o ID e o nome do jogador.
// Esta estrutura é usada para fornecer uma visão simplificada dos jogadores inscritos em um torneio,
// permitindo que a API retorne informações relevantes sobre cada jogador de forma concisa.
// A estrutura JogadorDetalhes é essencial para o gerenciamento de inscrições em torneios,
// permitindo que a aplicação trate tanto jogadores individuais quanto duplas de forma consistente.
// Ela é usada para gerar respostas detalhadas de inscrição, incluindo informações sobre o jogador
// ou dupla inscrita, facilitando a visualização e o gerenciamento das inscrições pelos administradores
// e organizadores de torneios.
// A estrutura também é usada para validar os dados de entrada ao inscrever um jogador ou dupla
// em um torneio, garantindo que os campos obrigatórios sejam preenchidos corretamente.
// A estrutura JogadorDetalhes é usada para receber os dados de entrada da API
// ao inscrever um jogador ou dupla em um torneio.
// Ela garante que os campos obrigatórios sejam preenchidos corretamente e que o tipo de modalidade
// seja válido (simples ou duplas).
// Ela também é usada para gerar respostas detalhadas de inscrição, incluindo informações sobre o jogador
// ou dupla inscrita, facilitando a visualização e o gerenciamento das inscrições pelos administradores
// e organizadores de torneios.
//	@Description	JogadorDetalhes é uma estrutura que contém informações simplificadas de um jogador para a resposta da API.
//	@ID				JogadorDetalhes
//	@Name			JogadorDetalhes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	JogadorDetalhes
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/jogadores_detalhes [get]
//	@Router			/jogadores_detalhes/{id} [get]
//	@Router			/jogadores_detalhes [post]
//	@Router			/jogadores_detalhes/{id} [put]
//	@Router			/jogadores_detalhes/{id} [delete]
//	@Tags			JogadoresDetalhes
//	@Param			id					path	int				true	"ID do JogadorDetalhes"
//	@Param			jogador_detalhes	body	JogadorDetalhes	true	"Dados do JogadorDetalhes"
//	@Param			nome				query	string			false	"Nome do Jogador"
//	@Param			page				query	int				false	"Número da página para paginação"
//	@Param			limit				query	int				false	"Número de itens por página para paginação"
//	@Param			sort				query	string			false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search				query	string			false	"Termo de busca para filtrar jogadores pelo nome"
//	@Security		BearerAuth
//	@Security		JWTAuth
type JogadorDetalhes struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

// DuplaDetalhes contém informações de uma dupla, incluindo seus jogadores.
// Ele inclui o ID da dupla, o nome da dupla (opcional),
// e detalhes dos jogadores A e B (opcionais).
// Esta estrutura é usada para fornecer uma visão detalhada das duplas inscritas em um torneio,
// permitindo que a API retorne informações relevantes sobre cada dupla de forma concisa.
// A estrutura DuplaDetalhes é essencial para o gerenciamento de inscrições em torneios,
// permitindo que a aplicação trate tanto jogadores individuais quanto duplas de forma consistente.
// Ela é usada para gerar respostas detalhadas de inscrição, incluindo informações sobre o jogador
// ou dupla inscrita, facilitando a visualização e o gerenciamento das inscrições pelos administradores
// e organizadores de torneios.
// A estrutura também é usada para validar os dados de entrada ao inscrever um jogador ou dupla
// em um torneio, garantindo que os campos obrigatórios sejam preenchidos corretamente.
// A estrutura DuplaDetalhes é usada para receber os dados de entrada da API
// ao inscrever um jogador ou dupla em um torneio.
// Ela garante que os campos obrigatórios sejam preenchidos corretamente e que o tipo de modalidade
// seja válido (simples ou duplas).
// Ela também é usada para gerar respostas detalhadas de inscrição, incluindo informações sobre o jogador
// ou dupla inscrita, facilitando a visualização e o gerenciamento das inscrições pelos administradores
// e organizadores de torneios.
//	@Description	DuplaDetalhes é uma estrutura que contém informações de uma dupla, incluindo seus jogadores.
//	@ID				DuplaDetalhes
//	@Name			DuplaDetalhes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	DuplaDetalhes
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/duplas_detalhes [get]
//	@Router			/duplas_detalhes/{id} [get]
//	@Router			/duplas_detalhes [post]
//	@Router			/duplas_detalhes/{id} [put]
//	@Router			/duplas_detalhes/{id} [delete]
//	@Tags			DuplasDetalhes
//	@Param			id				path	int				true	"ID da DuplaDetalhes"
//	@Param			dupla_detalhes	body	DuplaDetalhes	true	"Dados da DuplaDetalhes"
//	@Param			nome_dupla		query	string			false	"Nome da Dupla"
//	@Param			page			query	int				false	"Número da página para paginação"
//	@Param			limit			query	int				false	"Número de itens por página para paginação"
//	@Param			sort			query	string			false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search			query	string			false	"Termo de busca para filtrar duplas pelo nome"
//	@Security		BearerAuth
//	@Security		JWTAuth
type DuplaDetalhes struct {
	ID        int              `json:"id"`
	NomeDupla *string          `json:"nome_dupla,omitempty"`
	JogadorA  *JogadorDetalhes `json:"jogador_a,omitempty"`
	JogadorB  *JogadorDetalhes `json:"jogador_b,omitempty"`
}

// InscricaoDetalhada é a struct de resposta para a lista de inscritos em um torneio.
// Ela inclui o ID da inscrição, o tipo de modalidade (simples ou duplas),
// e, dependendo do tipo, inclui detalhes do jogador ou da dupla.
// Esta estrutura é usada para fornecer uma visão detalhada das inscrições,
// permitindo que a API retorne informações relevantes sobre cada inscrição.
// A estrutura InscricaoDetalhada é essencial para o gerenciamento de inscrições em torneios,
// permitindo que a aplicação trate tanto jogadores individuais quanto duplas de forma consistente.
// Ela é usada para gerar respostas detalhadas de inscrição, incluindo informações sobre o jogador
// ou dupla inscrita, facilitando a visualização e o gerenciamento das inscrições pelos administradores
// e organizadores de torneios.
// A estrutura também é usada para validar os dados de entrada ao inscrever um jogador ou dupla
// em um torneio, garantindo que os campos obrigatórios sejam preenchidos corretamente.
// A estrutura InscricaoDetalhada é usada para receber os dados de entrada da API
// ao inscrever um jogador ou dupla em um torneio.
// Ela garante que os campos obrigatórios sejam preenchidos corretamente e que o tipo de modalidade
// seja válido (simples ou duplas).
// Ela também é usada para gerar respostas detalhadas de inscrição, incluindo informações sobre o jogador
// ou dupla inscrita, facilitando a visualização e o gerenciamento das inscrições pelos administradores
// e organizadores de torneios.
//	@Description	InscricaoDetalhada é uma estrutura que representa uma inscrição detalhada em um torneio.
//	@ID				InscricaoDetalhada
//	@Name			InscricaoDetalhada
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	InscricaoDetalhada
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/inscricoes_detalhadas [get]
//	@Router			/inscricoes_detalhadas/{id} [get]
//	@Router			/inscricoes_detalhadas [post]
//	@Router			/inscricoes_detalhadas/{id} [put]
//	@Router			/inscricoes_detalhadas/{id} [delete]
//	@Tags			InscricoesDetalhadas
//	@Param			id					path	int					true	"ID da Inscrição Detalhada"
//	@Param			inscricao_detalhada	body	InscricaoDetalhada	true	"Dados da Inscrição Detalhada"
//	@Param			inscricao_id		query	int					false	"ID da Inscrição"
//	@Param			tipo_modalidade		query	string				false	"Tipo de Modalidade (simples ou duplas)"
//	@Param			page				query	int					false	"Número da página para paginação"
//	@Param			limit				query	int					false	"Número de itens por página para paginação"
//	@Param			sort				query	string				false	"Campo para ordenação, prefixado com '-' para ordem decrescente"
//	@Param			search				query	string				false	"Termo de busca para filtrar inscrições por jogador, dupla ou categoria"
//	@Security		BearerAuth
//	@Security		JWTAuth
type InscricaoDetalhada struct {
	InscricaoID    int              `json:"inscricao_id"`
	TipoModalidade string           `json:"tipo_modalidade"`
	Jogador        *JogadorDetalhes `json:"jogador,omitempty"` // Preenchido se for 'simples'
	Dupla          *DuplaDetalhes   `json:"dupla,omitempty"`   // Preenchido se for 'duplas'
}
