package handlers

import (
	"competitions/models"
	"competitions/repository"
	"competitions/validation"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// TorneioHandler encapsula a lógica para as rotas de torneios.
type TorneioHandler struct {
	repo repository.TorneioRepository
}

// NewTorneioHandler cria uma nova instância de TorneioHandler com o repositório fornecido.
// Ele é responsável por inicializar o handler com as dependências necessárias,
// como o repositório de torneios, que será usado para interagir com os dados
// relacionados aos torneios no banco de dados.
func NewTorneioHandler(repo repository.TorneioRepository) *TorneioHandler {
	return &TorneioHandler{repo: repo}
}

// CreateTorneio godoc
//
//	@Summary		Cria um novo torneio
//	@Description	Cria um novo torneio no sistema.
//	@Tags			Torneios
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		models.TorneioInput	true	"Dados do Torneio para Criação"
//	@Success		201		{object}	models.Torneio
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//
// swagger:route POST /torneios torneios CreateTorneio
// CreateTorneio cria um novo torneio no sistema.
// swagger:response CreateTorneio
// A resposta contém o objeto Torneio recém-criado, que inclui informações como ID,
// nome, data de início, data de término, local, tipo de modalidade e status do torneio.
// A resposta é retornada no formato JSON e pode ser usada para exibir os detalhes
// do torneio recém-criado. Em caso de erro ao criar o torneio, uma mensagem de erro será retornada
// com o status HTTP 500 (Internal Server Error).
// A validação dos dados de entrada é realizada para garantir que todas as informações necessárias
// estejam presentes e corretas. Se os dados forem inválidos, uma mensagem de erro será retornada
// com o status HTTP 400 (Bad Request), indicando que os dados fornecidos não são válidos.
// A validação inclui verificar campos obrigatórios, formatos de data e outros critérios
// específicos do torneio, como tipo de modalidade e status.
// Se a criação do torneio for bem-sucedida, a resposta incluirá o objeto Torneio com os dados
// do torneio recém-criado.
func (h *TorneioHandler) CreateTorneio(c *gin.Context) {
	var input models.TorneioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validação aprimorada dos dados de entrada
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	torneio, err := h.repo.Create(c.Request.Context(), input)
	if err != nil {
		log.Printf("Erro ao criar torneio: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno ao criar o torneio."})
		return
	}

	c.JSON(http.StatusCreated, torneio)
}

// GetTorneios godoc
//
//	@Summary		Lista todos os torneios
//	@Description	Retorna uma lista de todos os torneios registrados no sistema.
//	@Tags			Torneios
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.Torneio
//	@Failure		500	{object}	ErrorResponse
//
// swagger:route GET /torneios torneios GetTorneios
// GetTorneios retorna uma lista de todos os torneios registrados no sistema.
// swagger:response GetTorneios
// A resposta contém um array de objetos Torneio, cada um representando um torneio
// registrado no sistema. Cada objeto Torneio inclui informações como ID, nome, data de início,
// data de término, local, tipo de modalidade e status do torneio.
// A resposta é retornada no formato JSON e pode ser usada para exibir uma lista de
// torneios disponíveis para os usuários. Se não houver torneios registrados, a resposta
// será um array vazio. Em caso de erro ao buscar os torneios, uma mensagem de erro será retornada
// com o status HTTP 500 (Internal Server Error).
func (h *TorneioHandler) GetTorneios(c *gin.Context) {
	torneios, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		log.Printf("Erro ao buscar torneios: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os torneios."})
		return
	}

	c.JSON(http.StatusOK, torneios)
}

// GetTorneioByID godoc
//
//	@Summary		Busca um torneio por ID
//	@Description	Retorna um único torneio com base no ID fornecido.
//	@Tags			Torneios
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID do Torneio"
//	@Success		200	{object}	models.Torneio
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//
// swagger:route GET /torneios/{id} torneios GetTorneioByID
// GetTorneioByID busca um torneio específico pelo ID fornecido na URL.
// swagger:response GetTorneioByID
// A resposta contém o objeto Torneio correspondente ao ID fornecido, que inclui informações como
// ID, nome, data de início, data de término, local, tipo de modalidade e status do torneio.
// Se o torneio for encontrado, ele será retornado no formato JSON.
// Se o ID fornecido for inválido, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
func (h *TorneioHandler) GetTorneioByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	torneio, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Torneio não encontrado"})
			return
		}
		log.Printf("Erro ao buscar torneio por ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar o torneio."})
		return
	}

	c.JSON(http.StatusOK, torneio)
}

// UpdateTorneio godoc
//
//	@Summary		Atualiza um torneio existente
//	@Description	Atualiza os dados de um torneio existente com base no ID fornecido.
//	@Tags			Torneios
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"ID do Torneio"
//	@Param			input	body		models.TorneioInput	true	"Dados do Torneio para Atualização"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//
// swagger:route PUT /torneios/{id} torneios UpdateTorneio
// UpdateTorneio atualiza os dados de um torneio existente com base no ID fornecido.
// swagger:response UpdateTorneio
// A resposta contém uma mensagem de sucesso indicando que o torneio foi atualizado com sucesso.
// Se o ID fornecido for inválido, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
// Se o torneio não for encontrado, uma mensagem de erro será retornada com o status HTTP 404 (Not Found).
// Se ocorrer um erro ao atualizar o torneio, uma mensagem de erro será retornada com o status HTTP 500 (Internal Server Error).
// A validação dos dados de entrada é realizada para garantir que todas as informações necessárias
// estejam presentes e corretas. Se os dados forem inválidos, uma mensagem de erro será retornada
// com o status HTTP 400 (Bad Request), indicando que os dados fornecidos não são válidos.
// A validação inclui verificar campos obrigatórios, formatos de data e outros critérios
// específicos do torneio, como tipo de modalidade e status.
// Se a atualização do torneio for bem-sucedida, a resposta incluirá uma mensagem
func (h *TorneioHandler) UpdateTorneio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input models.TorneioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validação aprimorada dos dados de entrada
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	rowsAffected, err := h.repo.Update(c.Request.Context(), id, input)
	if err != nil {
		log.Printf("Erro ao atualizar torneio %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao atualizar o torneio."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Torneio não encontrado para atualizar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Torneio atualizado com sucesso"})
}

// DeleteTorneio godoc
//
//	@Summary		Deleta um torneio
//	@Description	Deleta um torneio do sistema com base no ID fornecido.
//	@Tags			Torneios
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID do Torneio"
//	@Success		200	{object}	SuccessResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//
// swagger:route DELETE /torneios/{id} torneios DeleteTorneio
// DeleteTorneio deleta um torneio do sistema com base no ID fornecido.
// swagger:response DeleteTorneio
// A resposta contém uma mensagem de sucesso indicando que o torneio foi deletado com sucesso.
// Se o ID fornecido for inválido, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
// Se o torneio não for encontrado, uma mensagem de erro será retornada com o status HTTP 404 (Not Found).
// Se ocorrer um erro ao deletar o torneio, uma mensagem de erro será retornada com o status HTTP 500 (Internal Server Error).
func (h *TorneioHandler) DeleteTorneio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	rowsAffected, err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		log.Printf("Erro ao deletar torneio %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao deletar o torneio."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Torneio não encontrado para deletar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Torneio deletado com sucesso"})
}

// InscreverJogador godoc
// godoc
//
//	@Summary		Inscreve jogador(es) em um torneio
//	@Description	Permite inscrever um jogador individual ou uma dupla em um torneio específico.
//	@Tags			Torneios
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"ID do Torneio"
//	@Param			input	body		models.JogadorTorneioInput	true	"Dados de Inscrição (Jogador ou Dupla)"
//	@Success		201		{object}	models.JogadorTorneio
//	@Failure		400		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//
// swagger:route POST /torneios/{id}/inscricoes torneios InscreverJogador
// InscreverJogador permite inscrever um jogador individual ou uma dupla em um torneio específico.
// swagger:response InscreverJogador
// A resposta contém o objeto JogadorTorneio, que representa a inscrição do jogador ou dupla no torneio.
// O objeto inclui informações como ID do jogador, ID da dupla (se aplicável),
// ID do torneio, tipo de modalidade (simples ou duplas) e status da inscrição.
// Se a inscrição for bem-sucedida, o status HTTP 201 (Created) será retornado com o objeto JogadorTorneio.
// Se os dados de entrada forem inválidos, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
// A validação dos dados de entrada é realizada para garantir que todas as informações necessárias
// estejam presentes e corretas. Se os dados forem inválidos, uma mensagem de erro será retornada
// com o status HTTP 400 (Bad Request), indicando que os dados fornecidos não são válidos.
// A validação inclui verificar campos obrigatórios, como ID do jogador ou dupla,
// tipo de modalidade e ID do torneio. Se o tipo de modalidade for "simples", a validação
// deve garantir que o ID da dupla não esteja presente.
func (h *TorneioHandler) InscreverJogador(c *gin.Context) {
	// 1. Obter o ID do torneio a partir da URL
	torneioID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do torneio inválido"})
		return
	}

	// 2. Fazer o bind do corpo da requisição JSON
	var input models.JogadorTorneioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido: " + err.Error()})
		return
	}

	// Atribui o ID do torneio da URL ao nosso struct de entrada.
	input.TorneioID = torneioID

	// 3. Validar os dados de entrada usando o validador customizado.
	// Esta abordagem centraliza a lógica de validação no modelo (struct tags),
	// tornando o handler mais limpo e a validação mais consistente e reutilizável.
	// A struct `models.JogadorTorneioInput` deve ter tags como:
	//   - TipoModalidade: `validate:"required,oneof=simples duplas"` (conforme schema do DB)
	//   - JogadorID:      `validate:"required_if=TipoModalidade simples,excluded_with=DuplaID"`
	//   - DuplaID:        `validate:"required_if=TipoModalidade duplas,excluded_with=JogadorID"`
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados de inscrição inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	// 4. Chamar o método do repositório para criar a inscrição
	jogadorInscrito, err := h.repo.InscreverJogador(c.Request.Context(), input.ToModel())
	if err != nil {
		// Tratamento de erros aprimorado para fornecer feedback mais útil ao cliente.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503": // foreign_key_violation
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido fornecido. O jogador, categoria, dupla ou torneio especificado não existe."})
				return
			case "23505": // unique_violation
				c.JSON(http.StatusConflict, gin.H{"error": "Este jogador ou dupla já está inscrito neste torneio/categoria."})
				return
			}
		}

		// Log do erro para depuração e retorno de um erro genérico para o cliente.
		log.Printf("Erro ao inscrever jogador no torneio %d: %v", torneioID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao processar a inscrição."})
		return
	}

	c.JSON(http.StatusCreated, jogadorInscrito)
}

// ListarInscricoes retorna uma lista detalhada de todos os jogadores e duplas
// inscritos em um torneio específico.
// godoc
//
//	@Router			/torneios/{id}/inscricoes [get]
//	@Description	Retorna uma lista detalhada de todos os jogadores e duplas inscritos em um torneio específico.
//	@ID				ListarInscricoes
//	@Tags			Torneios
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID do Torneio"
//	@Success		200	{array}	models.JogadorTorneio
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//
// swagger:route GET /torneios/{id}/inscricoes torneios ListarInscricoes
// ListarInscricoes retorna uma lista detalhada de todos os jogadores e duplas
// inscritos em um torneio específico.
// swagger:response ListarInscricoes
// A resposta contém um array de objetos JogadorTorneio, cada um representando uma inscrição
// de jogador ou dupla no torneio. Cada objeto JogadorTorneio inclui informações like
// ID do jogador, ID da dupla (se aplicável), ID do torneio,
// tipo de modalidade (simples ou duplas) e status da inscrição.
// Se o torneio for encontrado, a lista de inscrições será retornada no formato JSON.
func (h *TorneioHandler) ListarInscricoes(c *gin.Context) {
	torneioID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do torneio inválido"})
		return
	}

	inscricoes, err := h.repo.ListarInscricoesPorTorneio(c.Request.Context(), torneioID)
	if err != nil {
		log.Printf("Erro ao listar inscrições do torneio %d: %v", torneioID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar as inscrições."})
		return
	}

	c.JSON(http.StatusOK, inscricoes)
}
