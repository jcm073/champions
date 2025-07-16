# Diretrizes do Projeto - Gemini

Este documento estabelece as convenções, padrões e melhores práticas a serem seguidas no desenvolvimento deste projeto. O objetivo é manter um código limpo, consistente, legível e de fácil manutenção.

## 1. Princípios Fundamentais

- **Simplicidade**: Prefira código simples e direto. Evite complexidade e abstrações desnecessárias.
- **Legibilidade**: O código deve ser fácil de entender. Use nomes de variáveis e funções descritivos.
- **Documentação**: Todo código novo, especialmente endpoints de API, deve ser devidamente documentado.

## 2. Padrões de Código Go

### 2.1. Estrutura do Projeto

O projeto segue uma estrutura de pacotes baseada em responsabilidades. Mantenha a consistência:

- `routes/`: Define todas as rotas da API e as associa aos seus respectivos handlers.
- `handlers/`: Contém a lógica de negócio da aplicação. Cada handler é responsável por receber uma requisição HTTP, processá-la (usando repositórios) e retornar uma resposta.
- `repository/`: Camada de acesso a dados. Toda a comunicação com o banco de dados deve ser feita exclusivamente através dos repositórios. **Handlers não devem conter SQL.**
- `models/`: Define as estruturas de dados (structs) que representam as entidades do nosso domínio (ex: `Usuario`, `Torneio`).
- `config/`: Gerencia a configuração da aplicação, como a conexão com o banco de dados.
- `middleware/`: Contém os middlewares da aplicação, como autenticação e tratamento de erros.
- `validation/`: Lógica para validação de dados de entrada.

### 2.2. Nomenclatura

- **Variáveis**: Use `camelCase`.
- **Funções e Structs Exportadas**: Use `PascalCase`.
- **Pacotes**: Nomes curtos, em minúsculas e que descrevam seu propósito (ex: `handlers`, `models`).

### 2.3. Tratamento de Erros

- Erros devem ser tratados imediatamente. Verifique sempre o `err` retornado pelas funções.
- Retorne erros para a camada superior quando não puderem ser tratados no contexto atual.
- Nos handlers, use as funções de resposta padronizadas (ex: `responses.APIError`) para garantir que os erros da API sejam consistentes.

## 3. Banco de Dados (PostgreSQL)

### 3.1. Schema

- O arquivo `schema-pg.sql` é a fonte da verdade para a estrutura do banco de dados.
- Qualquer alteração no schema (criação/alteração de tabelas) deve ser refletida neste arquivo.

### 3.2. Padrão de Repositório

- Toda interação com o banco de dados **deve** passar pela camada de `repository`.
- Use transações (`tx`) para operações que envolvem múltiplas escritas no banco, garantindo a atomicidade.

## 4. API e Documentação (Swagger)

**Regra de Ouro**: Todo endpoint novo ou modificado **DEVE** ter sua documentação Swagger completa e atualizada.

### 4.1. Como Documentar um Endpoint

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

### 4.2. Gerando a Documentação

Após adicionar ou atualizar os comentários de documentação, regenere os arquivos do Swagger com o comando:

```bash
swag init
```
