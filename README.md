Projeto para gerenciar torneios de esportes em geral.

# Configurações de conexão com o banco de dados Postgresql usadas pela aplicação.
# As variáveis abaixo são carregadas automaticamente pelo Go usando o pacote godotenv.
# Certifique-se de que os valores estejam corretos para o seu ambiente.

# Arquivo .env
DB_USER=root
DB_PASSWORD=sua_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=campeonatos

# Chave secreta para assinar os tokens JWT. Deve ser uma string longa e aleatória.
JWT_SECRET=seu_segredo_super_secreto_aqui

# Air para auto reload de arquivos estaticos
go install github.com/air-verse/air@latest
air init
air -c .air.toml ou apenas digite air

# Exemplos de Inser de Usuario
{
  "name": "Inserir 10 usuários jogadores",
  "method": "POST",
  "url": "http://localhost:3000/usuarios",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": [
    {
      "tipo": "jogador",
      "nome": "Jogador 1",
      "username": "jogador1",
      "cpf": "000.000.001-91",
      "data_nascimento": "2000-01-01",
      "email": "jogador1@email.com",
      "password": "Senha123!",
      "telefone": "11999990001",
      "ativo": true
    },
    {
      "tipo": "jogador",
      "nome": "Jogador 2",
      "username": "jogador2",
      "cpf": "000.000.002-82",
      "data_nascimento": "2000-02-02",
      "email": "jogador2@email.com",
      "password": "Senha123!",
      "telefone": "11999990002",
      "ativo": true
    },
    {
      "tipo": "jogador",
      "nome": "Jogador 3",
      "username": "jogador3",
      "cpf": "000.000.003-73",
      "data_nascimento": "2000-03-03",
      "email": "jogador3@email.com",
      "password": "Senha123!",
      "telefone": "11999990003",
      "ativo": true
    },
    {
      "tipo": "jogador",
      "nome": "Jogador 4",
      "username": "jogador4",
      "cpf": "000.000.004-64",
      "data_nascimento": "2000-04-04",
      "email": "jogador4@email.com",
      "password": "Senha123!",
      "telefone": "11999990004",
      "ativo": true
    },
    {
      "tipo": "jogador",
      "nome": "Jogador 5",
      "username": "jogador5",
      "cpf": "000.000.005-55",
      "data_nascimento": "2000-05-05",
      "email": "jogador5@email.com",
      "password": "Senha123!",
      "telefone": "11999990005",
      "ativo": true
    },
    {
      "tipo": "jogador",
      "nome": "Jogador 6",
      "username": "jogador6",
      "cpf": "000.000.006-46",
      "data_nascimento": "2000-06-06",
      "email": "jogador6@email.com",
      "password": "Senha123!",
      "telefone": "11999990006",
      "ativo": true
    },
    {
      "tipo": "jogador",
      "nome": "Jogador 7",
      "username": "jogador7",
      "cpf": "000.000.007-37",
      "data_nascimento": "2000-07-07",
      "email": "jogador7@email.com",
      "password": "Senha123!",
      "telefone": "11999990007",
      "ativo": true
    },
    {
      "tipo": "jogador",
      "nome": "Jogador 8",
      "username": "jogador8",
      "cpf": "000.000.008-28",
      "data_nascimento": "2000-08-08",
      "email": "jogador8@email.com",
      "password": "Senha123!",
      "telefone": "11999990008",
      "ativo": true
    },
    {
      "tipo": "jogador",
      "nome": "Jogador 9",
      "username": "jogador9",
      "cpf": "000.000.009-19",
      "data_nascimento": "2000-09-09",
      "email": "jogador9@email.com",
      "password": "Senha123!",
      "telefone": "11999990009",
      "ativo": true
    },
    {
      "tipo": "jogador",
      "nome": "Jogador 10",
      "username": "jogador10",
      "cpf": "000.000.010-00",
      "data_nascimento": "2000-10-10",
      "email": "jogador10@email.com",
      "password": "Senha123!",
      "telefone": "11999990010",
      "ativo": true
    }
  ]
}

# Exemplos de Insert de um Torneio
{
  "id_esporte": 1,
  "nome": "Segundo Torneio Patricia Medrado de Pickball Masters",
  "descricao": "Torneio de Pickball Masters",
  "quantidade_quadras": 3
}