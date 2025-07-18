# Diretrizes do Projeto - Gemini

Este documento estabelece as convenções, padrões e melhores práticas a serem seguidas no desenvolvimento deste projeto. O objetivo é manter um código limpo, consistente, legível e de fácil manutenção.

## 1. Visão Geral do Projeto

Esta é uma API REST para gerenciar torneios esportivos e gerar estatiscas atráves do registro de resultados utilizando dados do rank e rating de jogadores.

Esta API REST é desenvolvida em Go versão 1.24 ou superior e utiliza o framework Gin 
A API Rest inclui todo o CRUD necessário para cadastro de usuários/jogadores, clubes, arenas, esportes e com principal objetivo realizar a gestão de torneios tais como criação dos grupos, chaveamentos, emparceiramento até definir quem são os três primeiros colocados.

Esta API REST trabalhará com dois modelos possivéis de campeonato e o criador do torneio deverá optar por um destes modelos no momento da criação do torneio.
1-) Modelo de seeds(sementes) semelhantes aos modelos de torneios de tenis e copa do mundo de futebol.
2-) Modelo de ranking um contra todos. 
Os vencedores serão os três com a maior quantidade de partidas ganhas. O criterio de desempate será o confronto direto entre os jogadores e o segundo criterio de desempate será a quantidade pontos conquistado ao longo do torneio.
Neste modelo deve ser levado em conta a categoria e o rating de cada jogador para a criação dos grupos.


## 2. Configuração do Ambiente de Desenvolvimento

### 2.1. Pré-requisitos

- Go (versão 1.2x ou superior)
- Gin (versão 1.8x ou superior)
- PostgreSQL
- Docker (Opcional, para rodar o banco de dados)
- `swag` CLI (para documentação da API)

### 2.2. Instalação

1.  **Clonar o repositório:**
    ```bash
    git clone [https://github.com/jcm073/champions.git]
    cd [champions]
    ```
2.  **Instalar dependências Go:**
    ```bash
    go mod tidy
    ```
3.  **Instalar o `swag` para documentação:**
    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```
4.  **Configurar o banco de dados:**
    - Certifique-se de que uma instância do PostgreSQL esteja em execução.
    - Crie um banco de dados para o projeto.
    - Execute o script `schema-pg.sql` para criar as tabelas:
      ```bash
      psql -U [seu-usuario] -d [nome-do-banco] -f schema-pg.sql
      ```
    - Configure as variáveis de ambiente ou o arquivo `config/database.go` com as credenciais do seu banco.

## 3. Comandos Essenciais

-   **Executar a aplicação (com live-reload):**
    *(Detectei o arquivo .air.toml, que sugere o uso do Air para live-reload)*
    ```bash
    air
    ```
-   **Executar a aplicação (padrão):**
    ```bash
    go run main.go
    ```
-   **Executar os testes:**
    ```bash
    go test ./...
    ```
-   **Gerar/Atualizar a documentação da API (Swagger):**
    ```bash
    swag init
    ```
-   **Formatar o código:**
    ```bash
    go fmt ./...
    ```

## 4. Princípios Fundamentais

-   **Simplicidade**: Prefira código simples e direto. Evite complexidade e abstrações desnecessárias.
-   **Legibilidade**: O código deve ser fácil de entender. Use nomes de variáveis e funções descritivos.
-   **Documentação**: Todo código novo, especialmente endpoints de API, deve ser devidamente documentado.

## 5. Padrões de Código Go

### 5.1. Estrutura do Projeto

O projeto segue uma estrutura de pacotes baseada em responsabilidades. Mantenha a consistência:

-   `routes/`: Define todas as rotas da API e as associa aos seus respectivos handlers.
-   `handlers/`: Contém a lógica de negócio da aplicação. Cada handler é responsável por receber uma requisição HTTP, processá-la (usando repositórios) e retornar uma resposta.
-   `repository/`: Camada de acesso a dados. Toda a comunicação com o banco de dados deve ser feita exclusivamente através dos repositórios. **Handlers não devem conter SQL.**
-   `models/`: Define as estruturas de dados (structs) que representam as entidades do nosso domínio (ex: `Usuario`, `Torneio`).
-   `config/`: Gerencia a configuração da aplicação, como a conexão com o banco de dados.
-   `middleware/`: Contém os middlewares da aplicação, como autenticação e tratamento de erros.
-   `validation/`: Lógica para validação de dados de entrada.

### 5.2. Nomenclatura

-   **Variáveis**: Use `camelCase`.
-   **Funções e Structs Exportadas**: Use `PascalCase`.
-   **Pacotes**: Nomes curtos, em minúsculas e que descrevam seu propósito (ex: `handlers`, `models`).

### 5.3. Tratamento de Erros

-   Erros devem ser tratados imediatamente. Verifique sempre o `err` retornado pelas funções.
-   Retorne erros para a camada superior quando não puderem ser tratados no contexto atual.
-   Nos handlers, use as funções de resposta padronizadas (ex: `responses.APIError`) para garantir que os erros da API sejam consistentes.

## 6. Banco de Dados (PostgreSQL)

### 6.1. Schema

-   O arquivo `schema-pg.sql` é a fonte da verdade para a estrutura do banco de dados.
-   Qualquer alteração no schema (criação/alteração de tabelas) deve ser refletida neste arquivo.

### 6.2. Padrão de Repositório

-   Toda interação com o banco de dados **deve** passar pela camada de `repository`.
-   Use transações (`tx`) para operações que envolvem múltiplas escritas no banco, garantindo a atomicidade.

## 7. API e Documentação (Swagger)

**Regra de Ouro**: Todo endpoint novo ou modificado **DEVE** ter sua documentação Swagger completa e atualizada.

### 7.1. Como Documentar um Endpoint

Adicione um bloco de comentários formatado logo acima da função do seu handler. Use as anotações do `swag` para descrever o endpoint.

**Exemplo de um Handler Documentado:**

```go
// @Summary      Cria um novo torneio
// @Description  Adiciona um novo torneio ao banco de dados com base nos dados fornecidos
// @Tags         Torneios
// @Accept       json
// @Produce      json
// @Param        torneio  body      models.TorneioRequest  true  "Dados do Torneio para Criação"
// @Success      201      {object}  models.Torneio
// @Failure      400      {object}  responses.ErrorResponse "Dados inválidos"
// @Failure      500      {object}  responses.ErrorResponse "Erro interno do servidor"
// @Router       /torneios [post]
func (h *TorneioHandler) CreateTorneio(c *gin.Context) {
    // ... lógica do handler
}
```