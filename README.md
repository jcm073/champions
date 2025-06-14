Projeto para gerenciar torneios de esportes em geral.

# Configurações de conexão com o banco de dados MySQL usadas pela aplicação.
# As variáveis abaixo são carregadas automaticamente pelo Go usando o pacote godotenv.
# Certifique-se de que os valores estejam corretos para o seu ambiente.

# Arquivo .env
DB_USER=root
DB_PASSWORD=sua_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=campeonatos

# Air para auto reload de arquivos estaticos
go install github.com/air-verse/air@latest
air init
air -c .air.toml ou apenas digite air