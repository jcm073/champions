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
)

// EsporteHandler contém a dependência do repositório.
type EsporteHandler struct {
	repo repository.EsporteRepository
}

// NewEsporteHandler cria um novo EsporteHandler.
// Ele recebe um repositório de esportes como dependência.
// Isso permite que o handler interaja com a camada de persistência de dados.
// A injeção de dependência facilita a testabilidade e a manutenção do código.
// O repositório deve implementar a interface EsporteRepository definida no pacote repository.
// Isso garante que o handler possa chamar os métodos necessários para criar, buscar, atualizar e deletar esportes.
// O handler é responsável por lidar com as requisições HTTP relacionadas aos esportes, como
// criação, listagem, busca por ID, atualização e deleção de esportes.
// Ele utiliza o repositório para realizar as operações de persistência, mantendo a lógica de negócios separada da lógica de apresentação.
// O handler também lida com a validação dos dados de entrada, retornando erros apropriados quando necessário.
// Isso garante que as operações realizadas no repositório sejam baseadas em dados válidos e consistentes.
// A criação do handler é feita através da função NewEsporteHandler, que recebe o repositório como parâmetro.
// Essa abordagem permite que o handler seja facilmente configurado com diferentes implementações de repositório, facilitando testes e extensões de funcionalidade.
// @handler EsporteHandler
// @param repo repository.EsporteRepository - Repositório de esportes utilizado pelo handler.
// @return *EsporteHandler - Retorna uma instância do EsporteHandler configurada com o repositório fornecido.
// @example
//
//	repo := repository.NewEsporteRepository(db)
//	esporteHandler := handlers.NewEsporteHandler(repo)
//	router := gin.Default()
//	router.POST("/esportes", esporteHandler.CreateEsporte)
//	router.GET("/esportes", esporteHandler.GetEsportes)
//	router.GET("/esportes/:id", esporteHandler.GetEsporteByID)
//	router.PUT("/esportes/:id", esporteHandler.UpdateEsporte)
//	router.DELETE("/esportes/:id", esporteHandler.DeleteEsporte)
//
// @description Cria um novo EsporteHandler com o repositório fornecido.
// O handler é responsável por gerenciar as operações relacionadas aos esportes, como criação, busca, atualização e deleção.
// Ele utiliza o repositório para interagir com a camada de persistência de dados,
// garantindo que as operações sejam realizadas de forma consistente e segura.
// A injeção de dependência permite que o handler seja facilmente testável e extensível,
// possibilitando a substituição do repositório por uma implementação mockada durante os testes.
// @see repository.EsporteRepository - Interface que define os métodos necessários para o repositório
// de esportes, como Create, FindAll, FindByID, Update e Delete.
// @see models.EsporteInput - Estrutura que representa os dados de entrada para a criação
// e atualização de um esporte, incluindo validação de dados.
// @see validation.TranslateError - Função que traduz erros de validação em mensagens amigáveis
// para o usuário, facilitando a compreensão dos problemas encontrados durante a validação dos dados de
// entrada.
// @see gin.Context - Contexto da requisição HTTP utilizado pelo Gin para manipular as requisições e respostas.
// Ele permite acessar os dados da requisição, manipular o fluxo de execução e enviar respostas HTTP.
// @see pgx.ErrNoRows - Erro retornado pelo driver pgx quando uma consulta não retorna nenhum resultado.
// Esse erro é tratado especificamente para retornar uma resposta adequada ao usuário quando um esporte não é encontrado.
// @see log.Printf - Função utilizada para registrar mensagens de log, incluindo erros e informações
// relevantes durante a execução do handler. Isso ajuda na depuração e monitoramento do sistema.
// @see http.StatusBadRequest - Código de status HTTP 400, utilizado para indicar que a requisição
// contém dados inválidos ou malformados. É retornado quando a validação dos dados de entrada falha
// ou quando o ID fornecido na URL não é um número válido.
// @see http.StatusCreated - Código de status HTTP 201, utilizado para indicar que um recurso foi criado com sucesso.
// É retornado quando um novo esporte é criado com sucesso no repositório.
// @see http.StatusOK - Código de status HTTP 200, utilizado para indicar que a requisição foi processada com sucesso.
// É retornado quando a lista de esportes é recuperada, quando um esporte é encontrado por ID ou quando um esporte é atualizado com sucesso.
// @see http.StatusNotFound - Código de status HTTP 404, utilizado para indicar que o recurso solicitado não foi encontrado.
// É retornado quando um esporte não é encontrado no repositório, seja ao buscar por ID, ao atualizar ou ao deletar.
// @see http.StatusInternalServerError - Código de status HTTP 500, utilizado para indicar que ocorreu um erro interno no servidor.
// É retornado quando há falhas ao interagir com o repositório, como erros de banco de dados ou problemas de conexão.
// @see gin.H - Tipo utilizado pelo Gin para representar um mapa de dados que será
// convertido em JSON na resposta HTTP. É utilizado para enviar mensagens de erro, sucesso e dados
// de esportes de forma estruturada e fácil de entender pelo cliente.
// @see strconv.Atoi - Função utilizada para converter uma string em um número inteiro.
// É utilizada para converter o ID do esporte fornecido na URL de uma string para um inteiro,
// permitindo que o handler interaja corretamente com o repositório.
// @see c.ShouldBindJSON - Método do Gin utilizado para vincular os dados JSON da requisição
// a uma estrutura Go. Ele é utilizado para extrair os dados de entrada do corpo da requisição
// e preenchê-los na estrutura EsporteInput, que é utilizada para criar ou atualizar um esporte.
// @see c.JSON - Método do Gin utilizado para enviar uma resposta JSON ao cliente.
// Ele é utilizado para retornar os resultados das operações realizadas pelo handler, como a criação de
// um novo esporte, a lista de esportes, a busca por um esporte específico, a atualização de um esporte e a deleção de um esporte.
// Ele também é utilizado para retornar mensagens de erro quando ocorrem problemas durante o processamento das requisições.
// @see c.Param - Método do Gin utilizado para acessar os parâmetros da rota.
// Ele é utilizado para obter o ID do esporte a partir da URL da requisição,
// permitindo que o handler identifique qual esporte deve ser buscado,
// atualizado ou deletado. O ID é passado como um parâmetro na rota, e o handler o utiliza para realizar as operações correspondentes no repositório.
// @see c.Request.Context - Método do Gin utilizado para obter o contexto da requisição.
// Ele é utilizado para passar o contexto da requisição para o repositório, permitindo que as operações de banco de dados sejam executadas dentro
// do contexto da requisição atual. Isso é importante para garantir que as operações sejam canceláveis
// e que os recursos sejam liberados corretamente após o processamento da requisição.
func NewEsporteHandler(repo repository.EsporteRepository) *EsporteHandler {
	return &EsporteHandler{repo: repo}
}

// CreateEsporte é o handler para criar um novo esporte.
// Ele recebe os dados do esporte no corpo da requisição, valida esses dados e, se válidos,
// chama o repositório para criar o esporte no banco de dados.
// Se a validação falhar, retorna um erro 400 Bad Request com detalhes sobre os erros de validação.
// Se a criação for bem-sucedida, retorna o esporte criado com um status 201 Created.
// @handler CreateEsporte
// @param c *gin.Context - Contexto da requisição HTTP, utilizado para acessar os dados da requisição e enviar a resposta.
// @return void - Não retorna nenhum valor, mas envia uma resposta JSON ao cliente.
// @example
//
//	router.POST("/esportes", esporteHandler.CreateEsporte)
//
// @description Cria um novo esporte com os dados fornecidos no corpo da requisição.
// O handler valida os dados de entrada e, se válidos, chama o repositório
// para criar o esporte no banco de dados. Se a validação falhar, retorna um erro 400 Bad Request com detalhes sobre os erros de validação.
// Se a criação for bem-sucedida, retorna o esporte criado com um status 201 Created.
// @see models.EsporteInput - Estrutura que representa os dados de entrada para a criação de um esporte.
// Ela contém os campos necessários para criar um esporte, como nome, descrição e outros atributos relevantes.
// A estrutura também possui um método Validate que verifica se os dados estão corretos e completos.
// @see validation.TranslateError - Função que traduz erros de validação em mensagens amigáveis
// para o usuário, facilitando a compreensão dos problemas encontrados durante a validação dos dados de entrada.
// @see gin.Context - Contexto da requisição HTTP utilizado pelo Gin para manipular as requisições e respostas.
// Ele permite acessar os dados da requisição, manipular o fluxo de execução e enviar respostas HTTP.
// @see http.StatusBadRequest - Código de status HTTP 400, utilizado para indicar que a requisição
// contém dados inválidos ou malformados. É retornado quando a validação dos dados de entrada falha
// ou quando os dados fornecidos não atendem aos critérios de validação definidos na estrutura EsporteInput.
// @see http.StatusCreated - Código de status HTTP 201, utilizado para indicar que um recurso foi criado com sucesso.
// É retornado quando um novo esporte é criado com sucesso no repositório.
// @see gin.H - Tipo utilizado pelo Gin para representar um mapa de dados que será
// convertido em JSON na resposta HTTP. É utilizado para enviar mensagens de erro, sucesso e dados
// de esportes de forma estruturada e fácil de entender pelo cliente.
// @see c.ShouldBindJSON - Método do Gin utilizado para vincular os dados JSON da requisição
// a uma estrutura Go. Ele é utilizado para extrair os dados de entrada do corpo da requisição
// e preenchê-los na estrutura EsporteInput, que é utilizada
// para criar um novo esporte.
// @see c.JSON - Método do Gin utilizado para enviar uma resposta JSON ao cliente.
// Ele é utilizado para retornar o esporte criado com um status 201 Created ou para retornar erros
// de validação com um status 400 Bad Request.
// @see c.Request.Context - Método do Gin utilizado para obter o contexto da requisição.
// Ele é utilizado para passar o contexto da requisição para o repositório, permitindo que as operações de banco de dados sejam executadas dentro
// do contexto da requisição atual. Isso é importante para garantir que as operações sejam canceláveis
// e que os recursos sejam liberados corretamente após o processamento da requisição.
// @see repository.EsporteRepository - Interface que define os métodos necessários para o repositório
// de esportes, como Create, FindAll, FindByID, Update e Delete.
// O repositório é responsável por interagir com o banco de dados e realizar as operações
// de persistência dos dados de esportes. Ele deve implementar os métodos definidos na interface
// para que o handler possa chamar as operações correspondentes.
// O repositório deve ser injetado no handler através da função NewEsporteHandler,
// permitindo que o handler utilize os métodos do repositório para criar, buscar, atualizar e deletar esportes.
// @see log.Printf - Função utilizada para registrar mensagens de log, incluindo erros e informações
// relevantes durante a execução do handler. Isso ajuda na depuração e monitoramento do sistema.
// É utilizado para registrar erros ao criar um esporte, facilitando a identificação de problemas
// durante o processo de criação.
func (h *EsporteHandler) CreateEsporte(c *gin.Context) {
	var input models.EsporteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	esporte, err := h.repo.Create(c.Request.Context(), input)
	if err != nil {
		log.Printf("Erro ao criar esporte: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao criar o esporte."})
		return
	}

	c.JSON(http.StatusCreated, esporte)
}

// GetEsportes é o handler para listar todos os esportes.
// Ele chama o repositório para buscar todos os esportes no banco de dados.
// Se ocorrer um erro ao buscar os esportes, retorna um erro 500 Internal Server Error.
// Se a busca for bem-sucedida, retorna a lista de esportes com um status 200 OK.
// @handler GetEsportes
// @param c *gin.Context - Contexto da requisição HTTP, utilizado para acessar os dados da requisição e enviar a resposta.
// @return void - Não retorna nenhum valor, mas envia uma resposta JSON ao cliente.
// @example
//
//	router.GET("/esportes", esporteHandler.GetEsportes)
//
// @description Busca todos os esportes cadastrados no sistema.
// O handler chama o repositório para buscar todos os esportes no banco de dados.
// Se ocorrer um erro ao buscar os esportes, retorna um erro 500 Internal Server Error.
// Se a busca for bem-sucedida, retorna a lista de esportes com um status 200 OK.
func (h *EsporteHandler) GetEsportes(c *gin.Context) {
	esportes, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		log.Printf("Erro ao buscar esportes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os esportes."})
		return
	}

	c.JSON(http.StatusOK, esportes)
}

// GetEsporteByID é o handler para buscar um esporte específico pelo ID.
// Ele extrai o ID do esporte dos parâmetros da rota, valida se é um número válido,
// e chama o repositório para buscar o esporte correspondente no banco de dados.
// Se o ID for inválido, retorna um erro 400 Bad Request.
// @handler GetEsporteByID
// @param c *gin.Context - Contexto da requisição HTTP, utilizado para acessar os dados da requisição e enviar a resposta.
// @return void - Não retorna nenhum valor, mas envia uma resposta JSON ao cliente.
// @example
//
//	router.GET("/esportes/:id", esporteHandler.GetEsporteByID)
//
// @description Busca um esporte específico pelo ID fornecido na rota.
// O handler extrai o ID do esporte dos parâmetros da rota, valida se é um número válido,
// e chama o repositório para buscar o esporte correspondente no banco de dados.
// Se o ID for inválido, retorna um erro 400 Bad Request.
// Se o esporte não for encontrado, retorna um erro 404 Not Found.
// Se ocorrer um erro ao buscar o esporte, retorna um erro 500 Internal Server Error.
// Se a busca for bem-sucedida, retorna o esporte encontrado com um status 200
func (h *EsporteHandler) GetEsporteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	esporte, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Esporte não encontrado"})
			return
		}
		log.Printf("Erro ao buscar esporte por ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar o esporte."})
		return
	}

	c.JSON(http.StatusOK, esporte)
}

// UpdateEsporte é o handler para atualizar um esporte existente.
// Ele extrai o ID do esporte dos parâmetros da rota, valida se é um número válido,
// e recebe os novos dados do esporte no corpo da requisição.
// Se os dados forem válidos, chama o repositório para atualizar o esporte no banco de dados.
// Se o ID for inválido, retorna um erro 400 Bad Request.
// @handler UpdateEsporte
// @param c *gin.Context - Contexto da requisição HTTP, utilizado para acessar os dados da requisição e enviar a resposta.
// @return void - Não retorna nenhum valor, mas envia uma resposta JSON ao cliente.
// @example
//
//	router.PUT("/esportes/:id", esporteHandler.UpdateEsporte)
//
// @description Atualiza um esporte existente com os novos dados fornecidos no corpo da requisição.
// O handler extrai o ID do esporte dos parâmetros da rota, valida se é um número válido,
// e recebe os novos dados do esporte no corpo da requisição.
// Se os dados forem válidos, chama o repositório para atualizar o esporte no banco de dados.
// Se o ID for inválido, retorna um erro 400 Bad Request.
// Se o esporte não for encontrado, retorna um erro 404 Not Found.
// Se ocorrer um erro ao atualizar o esporte, retorna um erro 500 Internal Server Error.
// Se a atualização for bem-sucedida, retorna uma mensagem de sucesso com um status
func (h *EsporteHandler) UpdateEsporte(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input models.EsporteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	rowsAffected, err := h.repo.Update(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao atualizar o esporte."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Esporte não encontrado para atualizar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Esporte atualizado com sucesso"})
}

// DeleteEsporte é o handler para deletar um esporte existente.
// Ele extrai o ID do esporte dos parâmetros da rota, valida se é um número válido,
// e chama o repositório para deletar o esporte correspondente no banco de dados.
// Se o ID for inválido, retorna um erro 400 Bad Request.
// @handler DeleteEsporte
// @param c *gin.Context - Contexto da requisição HTTP, utilizado para acessar os dados da requisição e enviar a resposta.
// @return void - Não retorna nenhum valor, mas envia uma resposta JSON ao cliente.
// @example
//
//	router.DELETE("/esportes/:id", esporteHandler.DeleteEsporte)
//
// @description Deleta um esporte existente pelo ID fornecido na rota.
// O handler extrai o ID do esporte dos parâmetros da rota, valida se é um número válido,
// e chama o repositório para deletar o esporte correspondente no banco de dados.
// Se o ID for inválido, retorna um erro 400 Bad Request.
// Se o esporte não for encontrado, retorna um erro 404 Not Found.
// Se ocorrer um erro ao deletar o esporte, retorna um erro 500 Internal Server Error.
// Se a deleção for bem-sucedida, retorna uma mensagem de sucesso com um status
func (h *EsporteHandler) DeleteEsporte(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	rowsAffected, err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		log.Printf("Erro ao deletar esporte %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao deletar o esporte."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Esporte não encontrado para deletar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Esporte deletado com sucesso"})
}
