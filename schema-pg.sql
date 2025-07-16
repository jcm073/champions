-- =================================================================================================
-- CRIAÇÃO DO BANCO DE DADOS (EXECUTAR SEPARADAMENTE ANTES DE RODAR O RESTANTE DO SCRIPT)
-- =================================================================================================
-- Certifique-se de que o locale 'pt_BR.UTF-8' esteja instalado no seu sistema operacional, se necessário:
--   sudo locale-gen pt_BR.UTF-8
--
-- Opção 1: Usando o comando `createdb` no terminal (recomendado):
--   sudo -u postgres createdb campeonatos --owner=postgres --encoding=UTF8 --locale=pt_BR.UTF-8 --template=template0
--
-- Opção 2: Usando um comando SQL no psql (conectado como superusuário ao BD 'postgres' ou 'template1'):
--   CREATE DATABASE campeonatos WITH OWNER postgres ENCODING 'UTF8' LC_COLLATE='pt_BR.UTF-8' LC_CTYPE='pt_BR.UTF-8' TEMPLATE template0;
--
-- Para listar bancos de dados existentes:
-- \l
--
-- =================================================================================================
-- EXECUÇÃO DO SCRIPT DE CRIAÇÃO DO ESQUEMA
-- =================================================================================================
-- APÓS A CRIAÇÃO DO BANCO, CONECTE-SE A ELE ANTES DE EXECUTAR O RESTANTE DESTE SCRIPT:
-- \c campeonatos
--
-- Para executar este script inteiro de uma vez (após criar e conectar ao banco):
--   psql -U seu_usuario -d campeonatos -f /home/jc/codigos/competitions/schema-pg.sql
--
-- Para executar em modo transacional (recomendado para garantir atomicidade, exceto para comandos que não podem estar em transações):
--   psql -U seu_usuario -d campeonatos --single-transaction -f /home/jc/codigos/competitions/schema-pg.sql
-- =================================================================================================

-- SEÇÃO 1: DEFINIÇÃO DE TIPOS ENUMERADOS (ENUMS)
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
        CREATE TYPE esporte_enum AS ENUM ('Beach Tenis','Tenis de Mesa','Tenis','Pickleball', 'Squash', 'Badminton', 'Padel');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tipo_modalidade_enum') THEN
        CREATE TYPE tipo_modalidade_enum AS ENUM ('simples','duplas');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'situacao_enum') THEN
        CREATE TYPE situacao_enum AS ENUM ('aguardando','em andamento','encerrado');
    END IF;
END$$;

-- SEÇÃO 2: TABELA DE USUÁRIOS
CREATE TABLE IF NOT EXISTS usuarios (
  id SERIAL PRIMARY KEY,
  tipo tipo_usuario NOT NULL,
  nome VARCHAR(100) NOT NULL,
  username VARCHAR(100) NOT NULL UNIQUE,
  cpf VARCHAR(14) NOT NULL UNIQUE,
  data_nascimento DATE NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  telefone VARCHAR(20) NOT NULL,
  instagram VARCHAR(50),
  criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ativo BOOLEAN NOT NULL DEFAULT TRUE
);

-- SEÇÃO 3: TABELA DE NÍVEIS
CREATE TABLE IF NOT EXISTS niveis (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(50) NOT NULL UNIQUE
);

-- SEÇÃO 4: TABELA DE TIPOS DE CATEGORIA
CREATE TABLE IF NOT EXISTS tipos_categoria (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(100) NOT NULL UNIQUE
);

-- SEÇÃO 5: TABELA DE CATEGORIAS
CREATE TABLE IF NOT EXISTS categorias (
  id SERIAL PRIMARY KEY,
  id_nivel INT NOT NULL REFERENCES niveis(id) ON DELETE CASCADE,
  descricao VARCHAR(100) NOT NULL,
  id_tipo_categoria INT NOT NULL REFERENCES tipos_categoria(id) ON DELETE CASCADE
);

-- SEÇÃO 6: TABELA DE ESPORTES
CREATE TABLE IF NOT EXISTS esportes (
  id SERIAL PRIMARY KEY,
  nome esporte_enum NOT NULL UNIQUE
);

-- SEÇÃO 7: TABELA DE PAÍSES
CREATE TABLE IF NOT EXISTS paises (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(100) NOT NULL UNIQUE
);

-- SEÇÃO 8: TABELA DE ESTADOS
CREATE TABLE IF NOT EXISTS estados (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(100) NOT NULL,
  sigla VARCHAR(2) NOT NULL UNIQUE,
  id_pais INT NOT NULL REFERENCES paises(id) ON DELETE CASCADE,
  UNIQUE (nome, id_pais) -- Garante que o nome do estado seja único dentro de um país
);

-- SEÇÃO 9: TABELA DE CIDADES
CREATE TABLE IF NOT EXISTS cidades (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(100) NOT NULL,
  id_estado INT NOT NULL REFERENCES estados(id) ON DELETE CASCADE,
  UNIQUE (nome, id_estado) -- Garante que o nome da cidade seja único dentro de um estado
);


-- 7. SCOUTS
CREATE TABLE IF NOT EXISTS scouts (
  id SERIAL PRIMARY KEY, -- Esta tabela parece ser para estatísticas de jogadores
  vitorias INT NOT NULL DEFAULT 0,
  derrotas INT NOT NULL DEFAULT 0,
  rating INT NOT NULL DEFAULT 0,
  titulos INT NOT NULL DEFAULT 0
);

-- SEÇÃO 10: JOGADORES
CREATE TABLE IF NOT EXISTS jogadores ( -- Representa um usuário que é um jogador
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

-- SEÇÃO 11: TABELA DE RELACIONAMENTO JOGADORES_ESPORTES (N:N)
CREATE TABLE IF NOT EXISTS jogadores_esportes (
  id_jogador INT NOT NULL REFERENCES jogadores(id) ON DELETE CASCADE,
  id_esporte INT NOT NULL REFERENCES esportes(id) ON DELETE CASCADE,
  PRIMARY KEY (id_jogador, id_esporte)
);

-- SEÇÃO 12: TABELA DE CLUBES
CREATE TABLE IF NOT EXISTS clubes (
  id SERIAL PRIMARY KEY,
  id_jogador_responsavel INT REFERENCES jogadores(id) ON DELETE SET NULL,
  nome VARCHAR(100) NOT NULL,
  telefone VARCHAR(20) NOT NULL,
  whatsapp VARCHAR(20),
  instagram VARCHAR(50),
  id_cidade INT NOT NULL REFERENCES cidades(id) ON DELETE CASCADE, -- Nova coluna
  id_estado INT NOT NULL REFERENCES estados(id) ON DELETE CASCADE, -- Nova coluna
  id_pais INT NOT NULL REFERENCES paises(id) ON DELETE CASCADE,
  quantidade INT NOT NULL DEFAULT 0,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- SEÇÃO 13: TABELA DE RELACIONAMENTO CLUBES_USUARIOS (N:N) - Jogadores em Clubes
CREATE TABLE IF NOT EXISTS clubes_usuarios (
  id_clube INT NOT NULL REFERENCES clubes(id) ON DELETE CASCADE,
  id_jogador INT NOT NULL REFERENCES jogadores(id) ON DELETE CASCADE,
  data_adesao TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id_clube, id_jogador)
);

-- SEÇÃO 14: TABELA DE TORNEIOS
CREATE TABLE IF NOT EXISTS torneios (
  id SERIAL PRIMARY KEY,
  id_esporte INT NOT NULL REFERENCES esportes(id) ON DELETE CASCADE,
  nome VARCHAR(100) NOT NULL,
  descricao TEXT,
  quantidade_quadras INT NOT NULL DEFAULT 1,
  criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Adicionada vírgula
  inicio TIMESTAMP NOT NULL, -- Renomeado de data_inicio
  fim TIMESTAMP NOT NULL,    -- Renomeado de data_fim
  id_cidade INT NOT NULL REFERENCES cidades(id) ON DELETE CASCADE, -- Nova coluna
  id_estado INT NOT NULL REFERENCES estados(id) ON DELETE CASCADE, -- Nova coluna
  id_pais INT NOT NULL REFERENCES paises(id) ON DELETE CASCADE,     -- Nova coluna
  ativo BOOLEAN NOT NULL DEFAULT TRUE
);

-- SEÇÃO 13: TABELA DE DUPLAS
CREATE TABLE IF NOT EXISTS duplas (
  id SERIAL PRIMARY KEY,
  id_jogador_a INT NOT NULL REFERENCES jogadores(id) ON DELETE CASCADE,
  id_jogador_b INT NOT NULL REFERENCES jogadores(id) ON DELETE CASCADE,
  nome_dupla VARCHAR(100),
  CONSTRAINT chk_jogadores_diferentes CHECK (id_jogador_a <> id_jogador_b)
  -- Optional: enforce order for simpler unique constraint on (A,B) vs (B,A)
  -- , CONSTRAINT chk_jogador_ordem CHECK (id_jogador_a < id_jogador_b) 
);
-- Optional: if chk_jogador_ordem is added
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_duplas_jogadores_unicos_ordenados ON duplas (id_jogador_a, id_jogador_b);

-- SEÇÃO 14: TABELA DE jogadores_torneios (jogadores associados aos torneios)
CREATE TABLE IF NOT EXISTS jogadores_torneios (
  id SERIAL PRIMARY KEY,
  id_torneio INT NOT NULL REFERENCES torneios(id) ON DELETE CASCADE,
  id_jogador INT REFERENCES jogadores(id) ON DELETE CASCADE,
  id_categoria INT NOT NULL REFERENCES categorias(id) ON DELETE CASCADE,
  id_dupla INT REFERENCES duplas(id) ON DELETE CASCADE,
  tipo_modalidade tipo_modalidade_enum NOT NULL DEFAULT 'simples',
  CONSTRAINT chk_jogador_torneio_modalidade_consistencia CHECK (
      (tipo_modalidade = 'simples' AND id_jogador IS NOT NULL AND id_dupla IS NULL) OR
      (tipo_modalidade = 'duplas' AND id_dupla IS NOT NULL AND id_jogador IS NULL)
  )
);

-- SEÇÃO 15: TABELA DE GRUPOS (de um torneio/categoria)
CREATE TABLE IF NOT EXISTS grupos (
  id SERIAL PRIMARY KEY,
  id_categoria INT NOT NULL REFERENCES categorias(id) ON DELETE CASCADE,
  nome VARCHAR(50) NOT NULL
);

-- SEÇÃO 16: TABELA DE RELACIONAMENTO GRUPO_JOGADORES_TORNEIOS (N:N)
CREATE TABLE IF NOT EXISTS grupo_jogadores_torneios (
  id_grupo INT NOT NULL REFERENCES grupos(id) ON DELETE CASCADE,
  id_jogador_torneio INT NOT NULL REFERENCES jogadores_torneios(id) ON DELETE CASCADE,
  PRIMARY KEY (id_grupo, id_jogador_torneio)
);

-- SEÇÃO 17: TABELA DE RODADAS (de um torneio/grupo)
CREATE TABLE IF NOT EXISTS rodadas (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(50) NOT NULL
);

-- 18. JOGOS
CREATE TABLE IF NOT EXISTS jogos (
  id SERIAL PRIMARY KEY, -- Jogos/Partidas
  id_torneio INT NOT NULL REFERENCES torneios(id) ON DELETE CASCADE,
  id_grupo INT NOT NULL REFERENCES grupos(id) ON DELETE CASCADE,
  id_rodada INT NOT NULL REFERENCES rodadas(id) ON DELETE CASCADE,
  id_jogador_torneio1 INT REFERENCES jogadores_torneios(id) ON DELETE SET NULL,
  id_jogador_torneio2 INT REFERENCES jogadores_torneios(id) ON DELETE SET NULL,
  id_dupla1 INT REFERENCES duplas(id) ON DELETE SET NULL,
  id_dupla2 INT REFERENCES duplas(id) ON DELETE SET NULL,
  id_jogador_vencedor INT REFERENCES jogadores(id) ON DELETE SET NULL,
  id_jogador_perdedor INT REFERENCES jogadores(id) ON DELETE SET NULL,
  id_dupla_vencedora INT REFERENCES duplas(id) ON DELETE SET NULL,
  id_dupla_perdedora INT REFERENCES duplas(id) ON DELETE SET NULL,
  tipo_modalidade tipo_modalidade_enum NOT NULL DEFAULT 'simples',
  data_hora TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  localizacao VARCHAR(100),
  situacao situacao_enum NOT NULL DEFAULT 'aguardando', -- Corrected default
  eh_final_campeonato BOOLEAN DEFAULT FALSE
);

-- SEÇÃO 19: TABELA DE SETS (scores por set)
CREATE TABLE IF NOT EXISTS sets (
  id SERIAL PRIMARY KEY,
  id_jogo INT NOT NULL REFERENCES jogos(id) ON DELETE CASCADE,
  set_numero INT NOT NULL,
  pontos_jogador1 INT NOT NULL,
  pontos_jogador2 INT NOT NULL,
  vencedor_set INT REFERENCES jogadores(id) ON DELETE SET NULL, -- ID do jogador que venceu o set
  UNIQUE (id_jogo, set_numero) -- Garante um score por set por jogo
);

-- SEÇÃO 20: CONSTRAINTS ADICIONAIS (ALTER TABLE)
-- Adicionar uma constraint para garantir a consistência dos dados de jogadores em torneios
-- Esta constraint garante que, para jogos 'simples', os campos de jogador do torneio sejam preenchidos e os de dupla sejam nulos,
-- e o oposto para jogos de 'duplas'.
ALTER TABLE jogos
    DROP CONSTRAINT IF EXISTS chk_jogo_participantes_modalidade; -- Remover se já existir uma com nome similar
ALTER TABLE jogos
    ADD CONSTRAINT chk_jogo_participantes_modalidade CHECK (
        -- This constraint checks the RESULT columns.
        (tipo_modalidade = 'simples' AND id_jogador_vencedor IS NOT NULL AND id_jogador_perdedor IS NOT NULL AND id_dupla_vencedora IS NULL AND id_dupla_perdedora IS NULL) 
        OR (tipo_modalidade = 'duplas' AND id_dupla_vencedora IS NOT NULL AND id_dupla_perdedora IS NOT NULL AND id_jogador_vencedor IS NULL AND id_jogador_perdedor IS NULL)
        -- Allow NULLs if results are not yet recorded
        OR (id_jogador_vencedor IS NULL AND id_jogador_perdedor IS NULL AND id_dupla_vencedora IS NULL AND id_dupla_perdedora IS NULL)
    );
ALTER TABLE jogos
    DROP CONSTRAINT IF EXISTS chk_jogo_definicao_participantes;
ALTER TABLE jogos
    DROP CONSTRAINT IF EXISTS chk_jogo_definicao_jogadores_torneios;
ALTER TABLE jogos
    ADD CONSTRAINT chk_jogo_definicao_jogadores_torneios CHECK (
        (tipo_modalidade = 'simples' AND id_jogador_torneio1 IS NOT NULL AND id_jogador_torneio2 IS NOT NULL AND id_dupla1 IS NULL AND id_dupla2 IS NULL) OR
        (tipo_modalidade = 'duplas' AND id_dupla1 IS NOT NULL AND id_dupla2 IS NOT NULL AND id_jogador_torneio1 IS NULL AND id_jogador_torneio2 IS NULL)
    );

ALTER TABLE duplas
    ADD CONSTRAINT chk_jogador_ordem CHECK (id_jogador_a < id_jogador_b);


-- SEÇÃO 21: ÍNDICES ÚTEIS
CREATE INDEX IF NOT EXISTS idx_jogadores_nome ON jogadores(nome);
CREATE INDEX IF NOT EXISTS idx_torneios_nome ON torneios(nome);
CREATE INDEX IF NOT EXISTS idx_jogos_data_hora ON jogos(data_hora);
CREATE INDEX IF NOT EXISTS idx_jogadores_torneios_torneio ON jogadores_torneios(id_torneio);
CREATE INDEX IF NOT EXISTS idx_jogadores_torneios_categoria ON jogadores_torneios(id_categoria);
CREATE INDEX IF NOT EXISTS idx_grupos_categoria ON grupos(id_categoria);
CREATE INDEX IF NOT EXISTS idx_jogos_torneio ON jogos(id_torneio);
CREATE INDEX IF NOT EXISTS idx_jogos_grupo ON jogos(id_grupo);
CREATE INDEX IF NOT EXISTS idx_jogos_rodada ON jogos(id_rodada);
CREATE INDEX IF NOT EXISTS idx_jogadores_torneios_jogador ON jogadores_torneios(id_jogador);
CREATE INDEX IF NOT EXISTS idx_jogadores_torneios_dupla ON jogadores_torneios(id_dupla);
CREATE INDEX IF NOT EXISTS idx_duplas_jogador_a ON duplas(id_jogador_a);
CREATE INDEX IF NOT EXISTS idx_duplas_jogador_b ON duplas(id_jogador_b);
CREATE UNIQUE INDEX IF NOT EXISTS idx_duplas_jogadores_unicos_ordenados ON duplas (id_jogador_a, id_jogador_b);
CREATE INDEX IF NOT EXISTS idx_placares_jogo ON placares(id_jogo);

-- SEÇÃO 22: FUNÇÕES E TRIGGERS

-- Função para criar um scout para um novo jogador se não for atribuído
CREATE OR REPLACE FUNCTION inserir_scout_para_novo_jogador()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.id_scout IS NULL THEN
        INSERT INTO scouts (vitorias, derrotas, rating, titulos)
        VALUES (0, 0, 0, 0)
        RETURNING id INTO NEW.id_scout;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_jogador_scout_before_insert ON jogadores;
CREATE TRIGGER trigger_jogador_scout_before_insert
BEFORE INSERT ON jogadores
FOR EACH ROW
EXECUTE FUNCTION inserir_scout_para_novo_jogador();

-- Função para atualizar a quantidade de membros em um clube
CREATE OR REPLACE FUNCTION atualizar_quantidade_membros_clube()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE clubes SET quantidade = quantidade + 1 WHERE id = NEW.id_clube;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE clubes SET quantidade = quantidade - 1 WHERE id = OLD.id_clube AND quantidade > 0;
    END IF;
    RETURN NULL; 
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_atualizar_membros_clube ON clubes_usuarios;
CREATE TRIGGER trigger_atualizar_membros_clube
AFTER INSERT OR DELETE ON clubes_usuarios
FOR EACH ROW
EXECUTE FUNCTION atualizar_quantidade_membros_clube();

-- Função para inserir jogador automaticamente ao criar usuário
CREATE OR REPLACE FUNCTION inserir_jogador_ao_criar_usuario()
RETURNS TRIGGER AS $$
BEGIN 
    IF NEW.tipo = 'jogador' THEN
        INSERT INTO jogadores (
            id_usuario, nome, cpf, data_nascimento, email, telefone, instagram, sexo, tipo, criado_em, ativo 
            -- id_scout será populado pelo trigger trigger_jogador_scout_before_insert
        ) VALUES (
            NEW.id, NEW.nome, NEW.cpf, NEW.data_nascimento, NEW.email, NEW.telefone, NEW.instagram, 
            'M'::sexo_enum, -- Defina um padrão ou colete essa informação no cadastro do usuário
            'destro'::tipo_mao_enum, -- Defina um padrão ou colete
            NOW(), NEW.ativo
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_usuario_jogador ON usuarios;
CREATE TRIGGER trigger_usuario_jogador
AFTER INSERT ON usuarios
FOR EACH ROW
EXECUTE FUNCTION inserir_jogador_ao_criar_usuario();

-- Trigger para sincronizar JOGADORES com USUARIOS (se mantiver colunas redundantes)
CREATE OR REPLACE FUNCTION sincronizar_jogador_com_usuario()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        -- Apenas atualiza se algum dos campos relevantes mudou
        IF OLD.nome IS DISTINCT FROM NEW.nome OR
           OLD.cpf IS DISTINCT FROM NEW.cpf OR
           OLD.data_nascimento IS DISTINCT FROM NEW.data_nascimento OR
           OLD.email IS DISTINCT FROM NEW.email OR
           OLD.telefone IS DISTINCT FROM NEW.telefone OR
           OLD.instagram IS DISTINCT FROM NEW.instagram THEN
           
            UPDATE jogadores
            SET nome = NEW.nome,
                cpf = NEW.cpf,
                data_nascimento = NEW.data_nascimento,
                email = NEW.email,
                telefone = NEW.telefone,
                instagram = NEW.instagram
            WHERE id_usuario = NEW.id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_sincronizar_jogador_com_usuario ON usuarios;
CREATE TRIGGER trigger_sincronizar_jogador_com_usuario
AFTER UPDATE OF nome, cpf, data_nascimento, email, telefone, instagram ON usuarios
FOR EACH ROW
EXECUTE FUNCTION sincronizar_jogador_com_usuario();


-- Função auxiliar para atualizar estatísticas de scout
CREATE OR REPLACE FUNCTION _aux_atualizar_estatisticas_scout(
    p_id_jogador INT,
    p_vitoria BOOLEAN,
    p_desfazer BOOLEAN, 
    p_incremento_rating_vitoria INT,
    p_decremento_rating_derrota INT,
    p_eh_final_campeonato BOOLEAN
)
RETURNS VOID AS $$
DECLARE
    v_operador_estatisticas INT := CASE WHEN p_desfazer THEN -1 ELSE 1 END;
    v_id_scout INT;
BEGIN
    IF p_id_jogador IS NULL THEN
        RETURN;
    END IF;

    SELECT id_scout INTO v_id_scout FROM jogadores WHERE id = p_id_jogador;

    IF v_id_scout IS NOT NULL THEN
        IF p_vitoria THEN
            UPDATE scouts
            SET vitorias = vitorias + (1 * v_operador_estatisticas),
                rating   = rating + (p_incremento_rating_vitoria * v_operador_estatisticas),
                titulos  = titulos + (CASE WHEN p_eh_final_campeonato THEN 1 ELSE 0 END * v_operador_estatisticas)
            WHERE id = v_id_scout;
        ELSE
            UPDATE scouts
            SET derrotas = derrotas + (1 * v_operador_estatisticas),
                -- Para rating em derrota, a lógica de "desfazer" uma subtração é uma adição.
                -- Se p_desfazer é TRUE, v_operador_estatisticas é -1.
                -- rating = rating - (val * op) => rating = rating - (val * -1) => rating = rating + val
                -- Se p_desfazer é FALSE, v_operador_estatisticas é 1.
                -- rating = rating - (val * op) => rating = rating - (val * 1)  => rating = rating - val
                rating   = rating - (p_decremento_rating_derrota * v_operador_estatisticas)
            WHERE id = v_id_scout;
        END IF;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Função principal do trigger para atualizar estatísticas de scout após um jogo
CREATE OR REPLACE FUNCTION atualizar_estatisticas_scout_jogo_completo()
RETURNS TRIGGER AS $$
DECLARE
    v_id_jogador_dupla_vencedora_a INT;
    v_id_jogador_dupla_vencedora_b INT;
    v_id_jogador_dupla_perdedora_a INT;
    v_id_jogador_dupla_perdedora_b INT;
    v_is_final_campeonato BOOLEAN; 
    v_rating_vitoria INT := 10; -- Valor de exemplo para incremento de rating em vitória
    v_rating_derrota INT := 5;  -- Valor de exemplo para decremento de rating em derrota
BEGIN
    v_is_final_campeonato := COALESCE(NEW.eh_final_campeonato, FALSE);

    IF TG_OP = 'INSERT' THEN
        IF NEW.eh_final_campeonato IS TRUE THEN
            v_is_final_campeonato := TRUE;
        END IF;

        IF NEW.tipo_modalidade = 'simples' THEN
            -- Atualiza scout do jogador vencedor
            PERFORM _aux_atualizar_estatisticas_scout(
                NEW.id_jogador_vencedor,
                TRUE, -- é vitória
                FALSE, -- não desfazer
                v_rating_vitoria,
                v_rating_derrota,
                v_is_final_campeonato
            );

            -- Atualiza scout do jogador perdedor
            PERFORM _aux_atualizar_estatisticas_scout(
                NEW.id_jogador_perdedor,
                FALSE, -- não é vitória (é derrota)
                FALSE, -- não desfazer
                v_rating_vitoria,
                v_rating_derrota,
                FALSE -- Títulos não são concedidos por derrota
            );

        ELSIF NEW.tipo_modalidade = 'duplas' THEN
            -- Obter jogadores da dupla vencedora
            SELECT id_jogador_a, id_jogador_b
            INTO v_id_jogador_dupla_vencedora_a, v_id_jogador_dupla_vencedora_b
            FROM duplas WHERE id = NEW.id_dupla_vencedora;

            -- Obter jogadores da dupla perdedora
            SELECT id_jogador_a, id_jogador_b
            INTO v_id_jogador_dupla_perdedora_a, v_id_jogador_dupla_perdedora_b
            FROM duplas WHERE id = NEW.id_dupla_perdedora;

            -- Atualizar scouts dos jogadores da dupla vencedora
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_vencedora_a, TRUE, FALSE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_vencedora_b, TRUE, FALSE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);

            -- Atualizar scouts dos jogadores da dupla perdedora
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_perdedora_a, FALSE, FALSE, v_rating_vitoria, v_rating_derrota, FALSE);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_perdedora_b, FALSE, FALSE, v_rating_vitoria, v_rating_derrota, FALSE);

        ELSE
            -- Modalidade não especificada ou desconhecida, pode ser útil logar um aviso
            RAISE WARNING 'Tipo de modalidade não especificado ou desconhecido para o jogo ID: %', NEW.id;
        END IF;

    ELSIF TG_OP = 'UPDATE' THEN        
        -- Desfazer estatísticas antigas (usando OLD)
        v_is_final_campeonato := COALESCE(OLD.eh_final_campeonato, FALSE);
        IF OLD.tipo_modalidade = 'simples' THEN
            PERFORM _aux_atualizar_estatisticas_scout(OLD.id_jogador_vencedor, TRUE, TRUE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(OLD.id_jogador_perdedor, FALSE, TRUE, v_rating_vitoria, v_rating_derrota, FALSE);
        ELSIF OLD.tipo_modalidade = 'duplas' THEN
            SELECT id_jogador_a, id_jogador_b INTO v_id_jogador_dupla_vencedora_a, v_id_jogador_dupla_vencedora_b FROM duplas WHERE id = OLD.id_dupla_vencedora;
            SELECT id_jogador_a, id_jogador_b INTO v_id_jogador_dupla_perdedora_a, v_id_jogador_dupla_perdedora_b FROM duplas WHERE id = OLD.id_dupla_perdedora;
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_vencedora_a, TRUE, TRUE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_vencedora_b, TRUE, TRUE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_perdedora_a, FALSE, TRUE, v_rating_vitoria, v_rating_derrota, FALSE);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_perdedora_b, FALSE, TRUE, v_rating_vitoria, v_rating_derrota, FALSE);
        END IF;

        -- Aplicar novas estatísticas (usando NEW, similar ao INSERT)
        v_is_final_campeonato := COALESCE(NEW.eh_final_campeonato, FALSE);
        IF NEW.tipo_modalidade = 'simples' THEN
            PERFORM _aux_atualizar_estatisticas_scout(NEW.id_jogador_vencedor, TRUE, FALSE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(NEW.id_jogador_perdedor, FALSE, FALSE, v_rating_vitoria, v_rating_derrota, FALSE);
        ELSIF NEW.tipo_modalidade = 'duplas' THEN
            SELECT id_jogador_a, id_jogador_b INTO v_id_jogador_dupla_vencedora_a, v_id_jogador_dupla_vencedora_b FROM duplas WHERE id = NEW.id_dupla_vencedora;
            SELECT id_jogador_a, id_jogador_b INTO v_id_jogador_dupla_perdedora_a, v_id_jogador_dupla_perdedora_b FROM duplas WHERE id = NEW.id_dupla_perdedora;
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_vencedora_a, TRUE, FALSE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_vencedora_b, TRUE, FALSE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_perdedora_a, FALSE, FALSE, v_rating_vitoria, v_rating_derrota, FALSE);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_perdedora_b, FALSE, FALSE, v_rating_vitoria, v_rating_derrota, FALSE);
        END IF;

    ELSIF TG_OP = 'DELETE' THEN
        -- Desfazer estatísticas do jogo excluído (usando OLD)
        v_is_final_campeonato := COALESCE(OLD.eh_final_campeonato, FALSE);
        IF OLD.tipo_modalidade = 'simples' THEN
            PERFORM _aux_atualizar_estatisticas_scout(OLD.id_jogador_vencedor, TRUE, TRUE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(OLD.id_jogador_perdedor, FALSE, TRUE, v_rating_vitoria, v_rating_derrota, FALSE);
        ELSIF OLD.tipo_modalidade = 'duplas' THEN
            SELECT id_jogador_a, id_jogador_b INTO v_id_jogador_dupla_vencedora_a, v_id_jogador_dupla_vencedora_b FROM duplas WHERE id = OLD.id_dupla_vencedora;
            SELECT id_jogador_a, id_jogador_b INTO v_id_jogador_dupla_perdedora_a, v_id_jogador_dupla_perdedora_b FROM duplas WHERE id = OLD.id_dupla_perdedora;
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_vencedora_a, TRUE, TRUE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_vencedora_b, TRUE, TRUE, v_rating_vitoria, v_rating_derrota, v_is_final_campeonato);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_perdedora_a, FALSE, TRUE, v_rating_vitoria, v_rating_derrota, FALSE);
            PERFORM _aux_atualizar_estatisticas_scout(v_id_jogador_dupla_perdedora_b, FALSE, TRUE, v_rating_vitoria, v_rating_derrota, FALSE);
        ELSE
            RAISE WARNING 'Tipo de modalidade não especificado ou desconhecido para o jogo ID: % ao tentar deletar.', OLD.id;
        END IF;
    END IF;
    RETURN NULL; -- Para AFTER triggers, o valor de retorno é ignorado
END;
$$ LANGUAGE plpgsql;

-- Trigger para atualizar estatísticas do scout após um jogo ser inserido,
-- utilizando a lógica mais completa.
-- Certifique-se que a tabela 'jogos' possui as colunas esperadas por
-- 'atualizar_estatisticas_scout_jogo_completo' (ex: id_jogador_vencedor, id_jogador_perdedor).
DROP TRIGGER IF EXISTS trigger_atualizar_estatisticas_scout_jogo ON jogos;
CREATE TRIGGER trigger_atualizar_estatisticas_scout_jogo
AFTER INSERT OR UPDATE OR DELETE ON jogos -- Adicionado UPDATE e DELETE para considerar esses casos
FOR EACH ROW
EXECUTE FUNCTION atualizar_estatisticas_scout_jogo_completo();

-- SEÇÃO 23: INSERÇÃO DE DADOS INICIAIS (SEEDING)
-- Esta seção é para popular tabelas com dados essenciais que são necessários para o funcionamento da aplicação.

-- Inserir os esportes definidos no ENUM na tabela de esportes.
-- O uso de 'ON CONFLICT (nome) DO NOTHING' garante que a execução repetida deste script não causará erros de duplicidade.
INSERT INTO esportes (nome) VALUES
('Badminton'),
('Beach Tenis'),
('Padel'),
('Pickleball'),
('Squash'),
('Tenis'),
('Tenis de Mesa')
ON CONFLICT (nome) DO NOTHING;
