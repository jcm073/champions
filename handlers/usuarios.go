package handlers

import (
	"competitions/middleware"
	"competitions/models"
	"competitions/repository"
	"competitions/validation"
	"errors"
	"log"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// UsuarioHandler holds the repository dependency.
type UsuarioHandler struct {
	repo repository.UsuarioRepository
}

// NewUsuarioHandler creates a new UsuarioHandler with the given repository.
func NewUsuarioHandler(repo repository.UsuarioRepository) *UsuarioHandler {
	return &UsuarioHandler{repo: repo}
}

// GetUsuarios godoc
//
//	@Summary		Lista todos os usuários
//	@Description	Retorna uma lista de todos os usuários registrados no sistema.
//	@Tags			Usuários
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.Usuario
//	@Failure		500	{object}	ErrorResponse
//
// Buscar todos os usuários
func (h *UsuarioHandler) GetUsuarios(c *gin.Context) {
	usuarios, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		log.Printf("Erro ao buscar usuários: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os usuários."})
		return
	}
	if len(usuarios) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhum usuário encontrado."})
		return
	}
	// O repositório já omite o campo de senha, então a resposta é segura.
	c.JSON(http.StatusOK, usuarios)
}

// GetUsuarioByID godoc
//
//	@Summary		Busca um usuário por ID
//	@Description	Retorna um único usuário com base no ID fornecido.
//	@Tags			Usuários
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID do Usuário"
//	@Success		200	{object}	models.Usuario
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//
// swagger:route GET /usuarios/{id} usuarios GetUsuarioByID
// GetUsuarioByID busca um usuário específico pelo ID fornecido na URL.
// swagger:response GetUsuarioByID
// A resposta contém o objeto Usuario correspondente ao ID fornecido, que inclui informações como
// ID, tipo, nome, username, CPF, data de nascimento, email, telefone, Instagram e status ativo.
// Se o usuário for encontrado, ele será retornado no formato JSON.
func (h *UsuarioHandler) GetUsuarioByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	usuario, err := h.repo.FindByID(c.Request.Context(), idInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		log.Printf("Erro ao buscar usuário por ID %d: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar o usuário."})
		return

	}
	c.JSON(http.StatusOK, usuario)
}

// CreateUsuario godoc
// @Summary		Cria um novo usuário
// @Description	Cria um novo usuário no sistema. Esta rota pode ser protegida por autenticação/autorização de admin.
// @Tags			Usuários
// @Accept			json
// @Produce		json
// @Param			input	body		models.UsuarioInput	true	"Dados do Usuário para Criação"
// @Success		201		{object}	models.Usuario
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Security		BearerAuth
// @Router			/usuarios [post]
// CreateUsuario cria um novo usuário no sistema.
// Esta rota deve ser protegida por autenticação/autorização de admin.
// A função recebe os dados do usuário no corpo da requisição, valida os dados, gera um hash para a senha e salva o usuário no banco de dados.
// Se a criação for bem-sucedida, retorna o usuário criado sem a senha.
// swagger:route POST /usuarios usuarios CreateUsuario
// CreateUsuario cria um novo usuário no sistema.
// swagger:response CreateUsuario
// A resposta contém o objeto Usuario criado, que inclui informações como ID, tipo, nome,
// username, CPF, data de nascimento, email, telefone, Instagram e status ativo.
// Se a criação for bem-sucedida, o campo Password será omitido para não expor informações sensíveis.
// Criar novo usuário
// NOTA: Esta função é similar à `Signup`, mas pode ser usada por um admin para criar usuários.
// A principal diferença é que esta rota deve ser protegida por autenticação/autorização de admin.
func (h *UsuarioHandler) CreateUsuario(c *gin.Context) {
	var input models.UsuarioInput
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

	// Hash da senha antes de salvar no banco
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao gerar hash da senha: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno."})
		return
	}

	usuario := models.Usuario{
		Tipo:           input.Tipo,
		Nome:           input.Nome,
		Username:       input.Username,
		CPF:            input.CPF,
		DataNascimento: input.DataNascimento,
		Email:          input.Email,
		Password:       string(hashedPassword),
		Telefone:       input.Telefone,
		Instagram:      input.Instagram,
		Ativo:          *input.Ativo,
	}

	err = h.repo.Create(c.Request.Context(), &usuario)
	if err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao criar o usuário."})
		return
	}
	// Limpa o campo Password para não expor a senha
	// Isso é importante para não expor informações sensíveis.
	usuario.Password = ""
	// O campo `Password` já tem `omitempty`, então ele não será serializado.
	// Retorne o usuário criado, mas sem a senha.
	c.JSON(http.StatusCreated, usuario)
}

// UpdateUsuario godoc
//
//	@Summary		Atualiza um usuário existente
//	@Description	Atualiza os dados de um usuário existente com base no ID fornecido.
//	@Tags			Usuários
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"ID do Usuário"
//	@Param			input	body		models.UpdateUsuarioInput	true	"Dados do Usuário para Atualização"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//
// swagger:route PUT /usuarios/{id} usuarios UpdateUsuario
// UpdateUsuario atualiza os dados de um usuário existente com base no ID fornecido.
// swagger:response UpdateUsuario
// A resposta contém uma mensagem de sucesso indicando que o usuário foi atualizado com sucesso.
// Se o ID fornecido for inválido, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
// Se o usuário não for encontrado, uma mensagem de erro será retornada com o status HTTP 404 (Not Found).
// Se ocorrer um erro ao atualizar o usuário, uma mensagem de erro será retornada com o status HTTP 500 (Internal Server Error).
func (h *UsuarioHandler) UpdateUsuario(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input models.UpdateUsuarioInput
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

	// 1. Buscar o usuário existente
	existingUser, err := h.repo.FindByID(c.Request.Context(), idInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado para atualizar"})
			return
		}
		log.Printf("Erro ao buscar usuário %d para atualização: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar o usuário."})
		return
	}

	// 2. Atualizar os campos do usuário existente com os dados da entrada
	existingUser.Tipo = input.Tipo
	existingUser.Nome = input.Nome
	existingUser.Username = input.Username
	existingUser.CPF = input.CPF
	existingUser.DataNascimento = input.DataNascimento
	existingUser.Email = input.Email
	existingUser.Telefone = input.Telefone
	existingUser.Instagram = input.Instagram
	existingUser.Ativo = *input.Ativo // Dereference the pointer

	// 3. Chamar o repositório para atualizar o usuário
	rowsAffected, err := h.repo.Update(c.Request.Context(), existingUser)
	if err != nil {
		log.Printf("Erro ao atualizar usuário %d: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao atualizar o usuário."})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado para atualizar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuário atualizado com sucesso"})
}

// DeleteUsuario godoc
//
//	@Summary		Deleta um usuário
//	@Description	Deleta um usuário do sistema com base no ID fornecido.
//	@Tags			Usuários
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID do Usuário"
//	@Success		200	{object}	SuccessResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//
// // swagger:route DELETE /usuarios/{id} usuarios DeleteUsuario
// DeleteUsuario deleta um usuário do sistema com base no ID fornecido.
// swagger:response DeleteUsuario
// A resposta contém uma mensagem de sucesso indicando que o usuário foi deletado com sucesso.
// Se o ID fornecido for inválido, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
// Se o usuário não for encontrado, uma mensagem de erro será retornada com o status HTTP 404 (Not Found).
// Se ocorrer um erro ao deletar o usuário, uma mensagem de erro será retornada com o status HTTP 500 (Internal Server Error).
func (h *UsuarioHandler) DeleteUsuario(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	rowsAffected, err := h.repo.Delete(c.Request.Context(), idInt)
	if err != nil {
		log.Printf("Erro ao deletar usuário %d: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao deletar o usuário."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado para deletar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário deletado com sucesso"})
}

// ChangePassword godoc
//
//	@Summary		Altera a senha do usuário
//	@Description	Altera a senha de um usuário específico, exigindo a senha antiga para verificação.
//	@Tags			Usuários
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"ID do Usuário"
//	@Param			input	body		models.ChangePasswordInput	true	"Senha Antiga e Nova Senha"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//
// swagger:route POST /usuarios/{id}/change-password usuarios ChangePassword
// ChangePassword altera a senha de um usuário específico, exigindo a senha antiga para verificação.
// swagger:response ChangePassword
// A resposta contém uma mensagem de sucesso indicando que a senha foi alterada com sucesso.
// Se o ID fornecido for inválido, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
// Se a senha antiga fornecida não corresponder à senha atual, uma mensagem de erro será retornada com o status HTTP 401 (Unauthorized).
// Se o usuário não for encontrado, uma mensagem de erro será retornada com o status HTTP404 (Not Found).
// Se ocorrer um erro ao alterar a senha, uma mensagem de erro será retornada com o status HTTP 500 (Internal Server Error).
// ChangePassword altera a senha de um usuário específico, exigindo a senha antiga para verificação.
// Esta rota deve ser protegida por autenticação/autorização de usuário.
// A função recebe o ID do usuário na URL e os dados da nova senha no corpo da requisição.
// Ela verifica se o ID do usuário na URL corresponde ao ID do usuário no token JWT, garantindo que o usuário só possa alterar
func (h *UsuarioHandler) ChangePassword(c *gin.Context) {
	// O ID do usuário na URL (:id) deve corresponder ao ID do usuário no token.
	// Isso impede que um usuário tente alterar a senha de outro.
	idParam := c.Param("id")
	idFromURL, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Error(&middleware.AppError{Code: http.StatusBadRequest, Message: "ID de usuário inválido na URL.", Err: err})
		return
	}

	// Extrai as claims do token JWT. O middleware já garantiu que o token é válido.
	claims := jwt.ExtractClaims(c)
	// O ID do usuário é armazenado como float64 após a decodificação do JSON do token.
	userIDClaim, ok := claims["user_id"].(float64)
	if !ok {
		c.Error(&middleware.AppError{Code: http.StatusUnauthorized, Message: "ID de usuário inválido ou ausente no token.", Err: errors.New("invalid token claims")})
		return
	}
	// Converte o ID para o tipo usado no modelo (uint).
	userID := uint(userIDClaim)

	var input models.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(&middleware.AppError{Code: http.StatusBadRequest, Message: "Corpo da requisição inválido.", Err: err})
		return
	}

	// Verificação de segurança crucial
	if uint(idFromURL) != userID {
		c.Error(&middleware.AppError{Code: http.StatusForbidden, Message: "Você não tem permissão para alterar a senha deste usuário.", Err: errors.New("permission denied")})
		return
	}

	if err := input.Validate(); err != nil {
		c.Error(&middleware.AppError{Code: http.StatusBadRequest, Message: "Dados inválidos.", Err: err})
		return
	}

	// 1. Buscar a senha hash atual do banco de dados
	// Usamos um método de repositório específico que busca a senha.
	// O método FindByID padrão omite a senha por segurança.
	user, err := h.repo.FindByIDForAuth(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.Error(&middleware.AppError{Code: http.StatusNotFound, Message: "Usuário não encontrado.", Err: err})
			return
		}
		c.Error(err) // Deixa o middleware tratar
		return
	}

	// 2. Comparar a senha antiga fornecida com o hash armazenado
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword))
	if err != nil {
		c.Error(&middleware.AppError{Code: http.StatusUnauthorized, Message: "Senha antiga incorreta.", Err: err})
		return
	}

	// 3. Gerar um novo hash para a nova senha
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Error(err)
		return
	}

	// 4. Atualizar a senha no banco de dados
	rowsAffected, err := h.repo.UpdatePassword(c.Request.Context(), userID, string(newHashedPassword))
	if err != nil {
		c.Error(err)
		return
	}

	if rowsAffected == 0 {
		c.Error(&middleware.AppError{Code: http.StatusNotFound, Message: "Usuário não encontrado para atualizar a senha.", Err: errors.New("user not found on update")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Senha atualizada com sucesso."})
}

// AssociateEsporte associa um jogador a um esporte
// godoc
//
//	@Summary		Associa um jogador a esportes
//	@Description	Associa um usuário (que deve ser um jogador) a um ou mais esportes.
//	@Tags			Usuários
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int								true	"ID do Usuário (Jogador)"
//	@Param			input	body		models.EsporteAssociationInput	true	"IDs dos Esportes para Associar"
//	@Success		201		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//
// swagger:route POST /usuarios/{id}/esportes usuarios AssociateEsporte
// AssociateEsporte associa um jogador a um ou mais esportes.
// swagger:response AssociateEsporte
// A resposta contém uma mensagem de sucesso indicando que o jogador foi associado ao(s) esporte(s) com sucesso.
// Se o ID do usuário fornecido for inválido, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
// Se o usuário não for encontrado, uma mensagem de erro será retornada com o status HTTP 404 (Not Found).
func (h *UsuarioHandler) AssociateEsporte(c *gin.Context) {
	// 1. Obter o ID do usuário da URL.
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	// 2. Obter os IDs dos esportes do corpo da requisição.
	var input models.EsporteAssociationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Validar a entrada.
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dados inválidos.",
			"errors":  validation.TranslateError(err),
		})
		return
	}

	// 4. Chamar o repositório para criar a associação.
	err = h.repo.AssociateEsporte(c.Request.Context(), userID, input.EsporteIDs)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrJogadorNaoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, repository.ErrEsporteInvalido):
			// Retorna 400 Bad Request, pois o cliente enviou dados inválidos (IDs de esporte que não existem).
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Printf("Erro ao associar esporte ao usuário %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao processar a associação."})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Jogador(es) associado(s) ao(s) esporte(s) com sucesso."})
}

// GetEsportesByUsuario retorna os esportes associados a um usuário específico
// godoc
//
//	@Summary		Busca os esportes associados a um usuário
//	@Description	Retorna uma lista de esportes associados a um usuário específico.
//	@Tags			Usuários
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID do Usuário (Jogador)"
//	@Success		200	{array}	models.Esporte
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//
// swagger:route GET /usuarios/{id}/esportes usuarios GetEsportesByUsuario
// GetEsportesByUsuario retorna os esportes associados a um usuário específico.
// swagger:response GetEsportesByUsuario
// A resposta contém uma lista de objetos Esporte associados ao usuário especificado pelo ID.
// Se o ID do usuário fornecido for inválido, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
// Se o usuário não for encontrado, uma mensagem de erro será retornada com o status HTTP 404 (Not Found).
func (h *UsuarioHandler) GetEsportesByUsuario(c *gin.Context) {
	// 1. Obter o ID do usuário da URL.
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	// 2. Chamar o repositório para buscar os esportes associados ao usuário.
	esportes, err := h.repo.GetEsportesByUsuario(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrJogadorNaoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		log.Printf("Erro ao buscar esportes para o usuário %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os esportes."})
		return
	}

	c.JSON(http.StatusOK, esportes)
}

// GetUsuariosByEsporte retorna os usuários associados a um esporte específico
// godoc
//
//	@Summary		Busca os usuários associados a um esporte
//	@Description	Retorna uma lista de usuários associados a um esporte específico.
//	@Tags			Usuários
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID do Esporte"
//	@Success		200	{array}	models.Usuario
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//
// swagger:route GET /esportes/{id}/usuarios usuarios GetUsuariosByEsporte
// GetUsuariosByEsporte retorna os usuários associados a um esporte específico.
// swagger:response GetUsuariosByEsporte
// A resposta contém uma lista de objetos Usuario associados ao esporte especificado pelo ID.
// Se o ID do esporte fornecido for inválido, uma mensagem de erro será retornada com o status HTTP 400 (Bad Request).
// Se o esporte não for encontrado, uma mensagem de erro será retornada com o status HTTP 404 (Not Found).
func (h *UsuarioHandler) GetUsuariosByEsporte(c *gin.Context) {
	// 1. Obter o ID do esporte da URL.
	esporteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de esporte inválido"})
		return
	}

	// 2. Chamar o repositório para buscar os usuários associados ao esporte.
	usuarios, err := h.repo.GetUsuariosByEsporte(c.Request.Context(), esporteID)
	if err != nil {
		if errors.Is(err, repository.ErrEsporteNaoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Esporte não encontrado"})
			return
		}
		log.Printf("Erro ao buscar usuários para o esporte %d: %v", esporteID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro ao buscar os usuários."})
		return
	}

	c.JSON(http.StatusOK, usuarios)
}
