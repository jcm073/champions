-- Criação do banco de dados (execute apenas uma vez, fora do psql conectado)
-- LC_COLLATE e LC_CTYPE devem ser compatíveis com o idioma do sistema operacional.
-- Se estiver usando o PostgreSQL no Linux, certifique-se de que o locale 'pt_BR.UTF-8' esteja instalado.
-- Caso não tenha o locale instalado, você pode criar um novo banco de dados com o locale correto:
-- sudo locale-gen pt_BR.UTF-8
-- Caso não tenha o banco de dados criado, execute o seguinte comando no terminal:
-- sudo -u postgres createdb campeonatos --encoding=UTF8 --locale=pt_BR.UTF-8 --template=template0
-- Caso queira criar o banco de dados diretamente no psql, execute:
-- CREATE DATABASE campeonatos WITH OWNER postgres ENCODING 'UTF8' LC_COLLATE='pt_BR.UTF-8' LC_CTYPE='pt_BR.UTF-8' TEMPLATE template0;

-- Listar bancos de dados
-- \l

-- Conecte-se ao banco antes de rodar o restante:
-- \c campeonatos

-- 1. ENUMS
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tipo_usuario') THEN
        CREATE TYPE tipo_usuario AS ENUM ('jogador', 'usuario', 'admin', 'gestor_clube', 'gestor_torneio');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sexo_enum') THEN
        CREATE TYPE sexo_enum AS ENUM ('M', 'F');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tipo_mao_enum') THEN
        CREATE TYPE tipo_mao_enum AS ENUM ('destro', 'canhoto');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'esporte_enum') THEN
        CREATE TYPE esporte_enum AS ENUM ('Beach Tenis','Tenis de Mesa','Tenis','Pickleball');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tipo_modalidade_enum') THEN
        CREATE TYPE tipo_modalidade_enum AS ENUM ('simples','duplas');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'situacao_enum') THEN
        CREATE TYPE situacao_enum AS ENUM ('aguardando','em andamento','encerrado');
    END IF;
END$$;

-- 2. USUARIOS
CREATE TABLE IF NOT EXISTS usuarios (
  id SERIAL PRIMARY KEY,
  tipo tipo_usuario NOT NULL,
  nome VARCHAR(100) NOT NULL,
  cpf VARCHAR(14) NOT NULL UNIQUE,
  data_nascimento DATE NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  telefone VARCHAR(20) NOT NULL,
  instagram VARCHAR(50),
  criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ativo BOOLEAN NOT NULL DEFAULT TRUE
);

-- 3. NIVEIS
CREATE TABLE IF NOT EXISTS niveis (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(50) NOT NULL UNIQUE
);

-- 4. TIPOS DE CATEGORIA
CREATE TABLE IF NOT EXISTS tipos_categoria (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(100) NOT NULL UNIQUE
);

-- 5. CATEGORIAS
CREATE TABLE IF NOT EXISTS categorias (
  id SERIAL PRIMARY KEY,
  id_nivel INT NOT NULL REFERENCES niveis(id) ON DELETE CASCADE,
  descricao VARCHAR(100) NOT NULL,
  id_tipo_categoria INT NOT NULL REFERENCES tipos_categoria(id) ON DELETE CASCADE
);

-- 6. ESPORTES
CREATE TABLE IF NOT EXISTS esportes (
  id SERIAL PRIMARY KEY,
  nome esporte_enum NOT NULL UNIQUE
);

-- 7. SCOUTS
CREATE TABLE IF NOT EXISTS scouts (
  id SERIAL PRIMARY KEY,
  vitorias INT NOT NULL DEFAULT 0,
  derrotas INT NOT NULL DEFAULT 0,
  pontos INT NOT NULL DEFAULT 0,
  titulos INT NOT NULL DEFAULT 0
);

-- 8. JOGADORES
CREATE TABLE IF NOT EXISTS jogadores (
  id SERIAL PRIMARY KEY,
  id_usuario INT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
  id_scout INT NOT NULL UNIQUE REFERENCES scouts(id) ON DELETE CASCADE,
  nome VARCHAR(100) NOT NULL,
  cpf VARCHAR(14) NOT NULL UNIQUE,
  data_nascimento DATE NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  telefone VARCHAR(20) NOT NULL,
  whatsapp VARCHAR(20),
  instagram VARCHAR(50),
  sexo sexo_enum NOT NULL,
  equipamento VARCHAR(255),
  tipo tipo_mao_enum NOT NULL,
  criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ativo BOOLEAN NOT NULL DEFAULT TRUE
);

-- 9. JOGADORES_ESPORTES (N:N)
CREATE TABLE IF NOT EXISTS jogadores_esportes (
  id_jogador INT NOT NULL REFERENCES jogadores(id) ON DELETE CASCADE,
  id_esporte INT NOT NULL REFERENCES esportes(id) ON DELETE CASCADE,
  PRIMARY KEY (id_jogador, id_esporte)
);

-- 10. CLUBES
CREATE TABLE IF NOT EXISTS clubes (
  id SERIAL PRIMARY KEY,
  id_jogador_responsavel INT REFERENCES jogadores(id) ON DELETE SET NULL,
  nome VARCHAR(100) NOT NULL,
  telefone VARCHAR(20) NOT NULL,
  whatsapp VARCHAR(20),
  instagram VARCHAR(50),
  estado VARCHAR(50) NOT NULL,
  cidade VARCHAR(100) NOT NULL,
  pais VARCHAR(50) NOT NULL DEFAULT 'Brasil',
  quantidade INT NOT NULL DEFAULT 0,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 11. CLUBES_USUARIOS (N:N)
CREATE TABLE IF NOT EXISTS clubes_usuarios (
  id_clube INT NOT NULL REFERENCES clubes(id) ON DELETE CASCADE,
  id_jogador INT NOT NULL REFERENCES jogadores(id) ON DELETE CASCADE,
  data_adesao TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id_clube, id_jogador)
);

-- 12. TORNEIOS
CREATE TABLE IF NOT EXISTS torneios (
  id SERIAL PRIMARY KEY,
  id_esporte INT NOT NULL REFERENCES esportes(id) ON DELETE CASCADE,
  nome VARCHAR(100) NOT NULL,
  descricao TEXT,
  quantidade_quadras INT NOT NULL DEFAULT 1,
  criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 13. DUPLAS
CREATE TABLE IF NOT EXISTS duplas (
  id SERIAL PRIMARY KEY,
  id_jogador_a INT NOT NULL REFERENCES jogadores(id) ON DELETE CASCADE,
  id_jogador_b INT NOT NULL REFERENCES jogadores(id) ON DELETE CASCADE,
  nome_dupla VARCHAR(100)
);

-- 14. PARTICIPANTES
CREATE TABLE IF NOT EXISTS participantes (
  id SERIAL PRIMARY KEY,
  id_torneio INT NOT NULL REFERENCES torneios(id) ON DELETE CASCADE,
  id_jogador INT REFERENCES jogadores(id) ON DELETE CASCADE,
  id_categoria INT NOT NULL REFERENCES categorias(id) ON DELETE CASCADE,
  id_dupla INT REFERENCES duplas(id) ON DELETE CASCADE,
  tipo_modalidade tipo_modalidade_enum NOT NULL DEFAULT 'simples'
);

-- 15. GRUPOS
CREATE TABLE IF NOT EXISTS grupos (
  id SERIAL PRIMARY KEY,
  id_categoria INT NOT NULL REFERENCES categorias(id) ON DELETE CASCADE,
  nome VARCHAR(50) NOT NULL
);

-- 16. GRUPO_PARTICIPANTES (N:N)
CREATE TABLE IF NOT EXISTS grupo_participantes (
  id_grupo INT NOT NULL REFERENCES grupos(id) ON DELETE CASCADE,
  id_participante INT NOT NULL REFERENCES participantes(id) ON DELETE CASCADE,
  PRIMARY KEY (id_grupo, id_participante)
);

-- 17. RODADAS
CREATE TABLE IF NOT EXISTS rodadas (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(50) NOT NULL
);

-- 18. JOGOS
CREATE TABLE IF NOT EXISTS jogos (
  id SERIAL PRIMARY KEY,
  id_torneio INT NOT NULL REFERENCES torneios(id) ON DELETE CASCADE,
  id_grupo INT NOT NULL REFERENCES grupos(id) ON DELETE CASCADE,
  id_rodada INT NOT NULL REFERENCES rodadas(id) ON DELETE CASCADE,
  id_participante1 INT REFERENCES participantes(id) ON DELETE SET NULL,
  id_participante2 INT REFERENCES participantes(id) ON DELETE SET NULL,
  id_dupla1 INT REFERENCES duplas(id) ON DELETE SET NULL,
  id_dupla2 INT REFERENCES duplas(id) ON DELETE SET NULL,
  tipo_modalidade tipo_modalidade_enum NOT NULL DEFAULT 'simples',
  data_hora TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  localizacao VARCHAR(100),
  situacao situacao_enum NOT NULL DEFAULT 'aguardando'
);

-- 19. PLACARES
CREATE TABLE IF NOT EXISTS placares (
  id SERIAL PRIMARY KEY,
  id_jogo INT REFERENCES jogos(id) ON DELETE SET NULL,
  resultado VARCHAR(100),
  id_jogador_vencedor INT REFERENCES jogadores(id) ON DELETE SET NULL,
  id_jogador_perdedor INT REFERENCES jogadores(id) ON DELETE SET NULL
);

-- 20. ÍNDICES ÚTEIS
CREATE INDEX IF NOT EXISTS idx_jogadores_nome ON jogadores(nome);
CREATE INDEX IF NOT EXISTS idx_torneios_nome ON torneios(nome);
CREATE INDEX IF NOT EXISTS idx_jogos_data_hora ON jogos(data_hora);
