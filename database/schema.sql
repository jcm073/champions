-- Cria o banco de dados
CREATE DATABASE IF NOT EXISTS campeonatos
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

-- Seleciona o banco
USE campeonatos;

-- Criação das tabelas

-- Tabela perfis
CREATE TABLE perfis (
  id INT AUTO_INCREMENT PRIMARY KEY,
  tipo ENUM('jogador', 'usuario', 'admin', 'gestor_clube', 'gestor_torneio') NOT NULL
) ENGINE=InnoDB;

-- Tabela usuarios
CREATE TABLE usuarios (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_perfil INT NOT NULL,
  nome VARCHAR(100) NOT NULL,
  cpf VARCHAR(14) NOT NULL, -- CONSIDERAÇÃO: Adicionar CONSTRAINT UNIQUE para cpf. Ex: ALTER TABLE usuarios ADD CONSTRAINT uc_usuario_cpf UNIQUE (cpf);
  datanascimento DATE NOT NULL,
  email VARCHAR(100) NOT NULL, -- CONSIDERAÇÃO: Adicionar CONSTRAINT UNIQUE para email. Ex: ALTER TABLE usuarios ADD CONSTRAINT uc_usuario_email UNIQUE (email);
  telefone VARCHAR(20) NOT NULL,
  instagram VARCHAR(50),
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  FOREIGN KEY (id_perfil) REFERENCES perfis(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
) ENGINE=InnoDB;

-- Tabela esportes
CREATE TABLE esportes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nome ENUM('beachtenis', 'tenisdemesa', 'tenis', 'pickleball') NOT NULL -- Renomeado 'name' para 'nome' e corrigido 'pickeball'
) ENGINE=InnoDB;

-- Tabela scouts (estatísticas gerais do jogador)
CREATE TABLE scouts ( -- Simplificada: removida referência circular complexa e colunas específicas de evento
  id INT AUTO_INCREMENT PRIMARY KEY,
  vitorias INT NOT NULL DEFAULT 0,
  derrotas INT NOT NULL DEFAULT 0,
  pontos INT NOT NULL DEFAULT 0,
  titulos INT NOT NULL DEFAULT 0
) ENGINE=InnoDB;

-- Tabela jogadores
CREATE TABLE jogadores (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_esportes INT NOT NULL, -- CONSIDERAÇÃO: Se um jogador puder praticar múltiplos esportes, uma tabela de junção (ex: jogadores_esportes(id_jogador, id_esporte)) seria mais flexível.
  id_scouts INT NOT NULL UNIQUE, -- Adicionado UNIQUE para relação 1:1 com scouts
  nome VARCHAR(100) NOT NULL,
  cpf VARCHAR(14) NOT NULL, -- CONSIDERAÇÃO: Adicionar CONSTRAINT UNIQUE para cpf. Ex: ALTER TABLE jogadores ADD CONSTRAINT uc_jogador_cpf UNIQUE (cpf);
  datanascimento DATE NOT NULL,
  email VARCHAR(100) NOT NULL, -- CONSIDERAÇÃO: Adicionar CONSTRAINT UNIQUE para email. Ex: ALTER TABLE jogadores ADD CONSTRAINT uc_jogador_email UNIQUE (email);
  telefone VARCHAR(20) NOT NULL,
  whatsapp VARCHAR(20), -- Corrigido 'whatsup'
  instagram VARCHAR(50),
  sexo ENUM('M', 'F') NOT NULL,
  equipamento VARCHAR(255),
  tipo ENUM('destro', 'canhoto') NOT NULL,
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id_esportes) REFERENCES esportes(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  FOREIGN KEY (id_scouts) REFERENCES scouts(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
) ENGINE=InnoDB;

-- Tabela clubes
CREATE TABLE clubes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_jogador_responsavel INT, -- Renomeado para clareza, se for um jogador específico
  nome VARCHAR(100) NOT NULL,
  responsavel VARCHAR(100) NOT NULL,
  telefone VARCHAR(20) NOT NULL,
  whatsapp VARCHAR(20), -- Corrigido 'whatsup'
  instagram VARCHAR(50),
  estado VARCHAR(50) NOT NULL,
  cidade VARCHAR(100) NOT NULL,
  pais VARCHAR(50) NOT NULL DEFAULT 'Brasil',
  quantidade INT NOT NULL DEFAULT 0, -- CONSIDERAÇÃO: Renomear para maior clareza (ex: quantidade_membros, capacidade_jogadores, quantidade_quadras_clube), dependendo do significado.
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id_jogador_responsavel) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
) ENGINE=InnoDB;

-- Tabela clubes_usuarios (tabela de junção)
CREATE TABLE clubes_usuarios (
  id_clubes INT NOT NULL,
  id_jogadores INT NOT NULL,
  descricao VARCHAR(255), -- Corrigido 'descrição'
  data_adesao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Corrigido 'datainserção' e renomeado
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id_clubes, id_jogadores),
  FOREIGN KEY (id_clubes) REFERENCES clubes(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  FOREIGN KEY (id_jogadores) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
) ENGINE=InnoDB;

-- Tabela niveis
CREATE TABLE niveis (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nome VARCHAR(50) NOT NULL UNIQUE -- Estrutura redesenhada para melhor representar níveis
) ENGINE=InnoDB;

-- Tabela categorias
CREATE TABLE categorias (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_nivel INT NOT NULL,
  descricao VARCHAR(100) NOT NULL,
  simplesmasculina BOOLEAN NOT NULL DEFAULT FALSE,
  simplesfeminina BOOLEAN NOT NULL DEFAULT FALSE,
  duplasmasculina BOOLEAN NOT NULL DEFAULT FALSE, -- Considerar normalização futura desta tabela
  duplasfeminina BOOLEAN NOT NULL DEFAULT FALSE,
  duplasmista BOOLEAN NOT NULL DEFAULT FALSE,
  40feminina BOOLEAN NOT NULL DEFAULT FALSE,
  40masculina BOOLEAN NOT NULL DEFAULT FALSE,
  50feminina BOOLEAN NOT NULL DEFAULT FALSE,
  50masculina BOOLEAN NOT NULL DEFAULT FALSE,
  60feminina BOOLEAN NOT NULL DEFAULT FALSE,
  60masculina BOOLEAN NOT NULL DEFAULT FALSE,
  veteranosfeminino BOOLEAN NOT NULL DEFAULT FALSE,
  veteranosmasculino BOOLEAN NOT NULL DEFAULT FALSE,
  FOREIGN KEY (id_nivel) REFERENCES niveis(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT -- CONSIDERAÇÃO DE NORMALIZAÇÃO: A estrutura atual com múltiplas colunas booleanas para tipos de categoria pode ser inflexível.
                       -- Sugestão: Uma tabela `tipos_categoria` (com `id`, `nome_tipo_categoria` como 'Simples Masculina', 'Duplas Feminina 40+', etc.)
                       -- e esta tabela `categorias` referenciaria `id_nivel` e `id_tipo_categoria`.
                       -- Alternativamente, colunas separadas para atributos como Gênero, Tipo de Jogo (Simples/Duplas), Faixa Etária.

) ENGINE=InnoDB;

-- Tabela torneios
CREATE TABLE torneios (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_esportes INT NOT NULL,
  nome VARCHAR(100) NOT NULL,
  descricao TEXT,
  quantidade_quadras INT NOT NULL DEFAULT 1, -- Corrigido 'quantidadequadras'
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id_esportes) REFERENCES esportes(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
) ENGINE=InnoDB;

-- Tabela participantes
CREATE TABLE participantes (
  -- CONSIDERAÇÃO GERAL DE DUPLAS: Esta tabela atualmente referencia id_jogadores.
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_torneios INT NOT NULL,
  id_jogadores INT NOT NULL,
  id_categorias INT NOT NULL,
  FOREIGN KEY (id_torneios) REFERENCES torneios(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  FOREIGN KEY (id_jogadores) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  FOREIGN KEY (id_categorias) REFERENCES categorias(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
  -- Se o sistema suportar duplas como entidades (ex: uma dupla se inscreve),
  -- esta tabela poderia ter `id_jogador` (para simples) e `id_dupla` (para duplas),
  -- ou um `id_participante_entidade` que aponta para uma tabela polimórfica ou duas colunas FK opcionais.
) ENGINE=InnoDB;


-- Tabela grupos
CREATE TABLE grupos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_categoria INT NOT NULL,
  nome VARCHAR(50) NOT NULL,
  id_jogador1 INT NOT NULL, -- CONSIDERAÇÃO: Número fixo de jogadores. Para flexibilidade, uma tabela de junção `grupo_participantes (id_grupo, id_participante)` seria melhor.
  id_jogador2 INT NOT NULL, -- `id_participante` poderia referenciar `participantes.id`.
  id_jogador3 INT,          -- CONSIDERAÇÃO DE DUPLAS: Como duplas são representadas aqui? Se id_jogador1 e id_jogador2 formam uma dupla, isso não está explícito.
  id_jogador4 INT,          -- Se grupos contêm duplas, a estrutura precisa refletir isso (ex: referenciar uma tabela `duplas` ou usar a tabela `participantes` adaptada).
  FOREIGN KEY (id_categoria) REFERENCES categorias(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  FOREIGN KEY (id_jogador1) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  FOREIGN KEY (id_jogador2) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  FOREIGN KEY (id_jogador3) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL,
  FOREIGN KEY (id_jogador4) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
) ENGINE=InnoDB;

-- Tabela rodadas
CREATE TABLE rodadas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nome VARCHAR(50) NOT NULL
) ENGINE=InnoDB;

-- Tabela placares (precisa ser criada antes de jogos devido à referência circular)
CREATE TABLE placares (
  id INT AUTO_INCREMENT PRIMARY KEY,
  pontuacao_jogador1 INT NOT NULL DEFAULT 0, -- Corrigido 'pontuação'
  pontuacao_jogador2 INT NOT NULL DEFAULT 0, -- Corrigido 'pontuação'
  resultado VARCHAR(100),
  id_jogador_vencedor INT, -- Alterado para ID, FK será adicionado após tabela jogadores
  id_jogador_perdedor INT  -- Alterado para ID, FK será adicionado após tabela jogadores
) ENGINE=InnoDB;

-- Tabela jogos
CREATE TABLE jogos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_torneios INT NOT NULL,
  id_grupos INT NOT NULL,
  id_rodadas INT NOT NULL,
  id_jogador1 INT NOT NULL, -- CONSIDERAÇÃO JOGOS DE DUPLAS: Adequado para simples. Para duplas, precisa adaptar.
  id_jogador2 INT NOT NULL, -- Opções: colunas adicionais (id_jogador3, id_jogador4), referenciar tabela `duplas` (id_dupla1, id_dupla2), ou referenciar `participantes` (id_participante1, id_participante2).
  id_placares INT NOT NULL,
  data_hora DATETIME NOT NULL,
  localizacao VARCHAR(100),
  situacao ENUM('aguardando', 'em andamento', 'encerrado') NOT NULL DEFAULT 'aguardando',
  FOREIGN KEY (id_torneios) REFERENCES torneios(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  FOREIGN KEY (id_grupos) REFERENCES grupos(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  FOREIGN KEY (id_rodadas) REFERENCES rodadas(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  FOREIGN KEY (id_jogador1) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  FOREIGN KEY (id_jogador2) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  FOREIGN KEY (id_placares) REFERENCES placares(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
) ENGINE=InnoDB;

-- Atualizar a tabela placares para adicionar a chave estrangeira
ALTER TABLE placares
  ADD COLUMN id_jogos INT,
  ADD FOREIGN KEY (id_jogos) REFERENCES jogos(id) ON UPDATE CASCADE ON DELETE SET NULL,
  ADD CONSTRAINT fk_placar_vencedor FOREIGN KEY (id_jogador_vencedor) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL,
  ADD CONSTRAINT fk_placar_perdedor FOREIGN KEY (id_jogador_perdedor) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL;

-- Adicionar aqui inserts para a tabela niveis, se desejar dados iniciais. Ex:
-- INSERT INTO niveis (nome) VALUES ('Iniciante'), ('Intermediário'), ('Intermediário Plus'), ('Avançado'), ('Profissional');

-- CONSIDERAÇÕES GERAIS ADICIONAIS DO ARQUIVO mudanças.txt:

-- 7. REPRESENTAÇÃO GERAL DE DUPLAS (Ponto Crucial):
--    Decidir como duplas serão gerenciadas. Se uma dupla é uma "equipe" que se inscreve,
--    joga em grupos e partidas, pode precisar ser uma entidade própria (ex: tabela `duplas`).
--    As tabelas `participantes`, `grupos`, `jogos` deveriam então referenciar `id_dupla`
--    para eventos de duplas e `id_jogador` para simples.
--    Uma coluna `tipo_modalidade` ENUM('simples', 'duplas') em `categorias` ou `torneios` pode ajudar.
--    Exemplo de tabela `duplas`:
--    CREATE TABLE duplas (
--      id INT AUTO_INCREMENT PRIMARY KEY,
--      id_jogador_a INT NOT NULL, FOREIGN KEY (id_jogador_a) REFERENCES jogadores(id),
--      id_jogador_b INT NOT NULL, FOREIGN KEY (id_jogador_b) REFERENCES jogadores(id),
--      nome_dupla VARCHAR(100)
--    );

-- 8. ÍNDICES:
--    Além de PKs e FKs, adicionar índices em colunas frequentemente usadas em `WHERE`, `JOIN`, `ORDER BY`, `GROUP BY`
--    para otimizar consultas (ex: `jogadores.nome`, `torneios.nome`, `jogos.data_hora`).

-- 9. REVISÃO DE `ON DELETE` E `ON UPDATE` NAS CHAVES ESTRANGEIRAS:
--    Revisar as ações. `RESTRICT` é seguro. `CASCADE` pode ser útil (ex: deletar clube -> remover de `clubes_usuarios`), mas usar com cautela. `SET NULL` é outra opção.