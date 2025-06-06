-- Cria o banco de dados
CREATE DATABASE IF NOT EXISTS campeonatos
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE campeonatos;

-- 1. PERFIS (ENUM direto em usuarios para simplificar)
-- Remover tabela perfis, usar ENUM em usuarios
-- 2. USUARIOS
CREATE TABLE usuarios (
  id INT AUTO_INCREMENT PRIMARY KEY,
  tipo ENUM('jogador', 'usuario', 'admin', 'gestor_clube', 'gestor_torneio') NOT NULL,
  nome VARCHAR(100) NOT NULL,
  cpf VARCHAR(14) NOT NULL UNIQUE,
  datanascimento DATE NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  telefone VARCHAR(20) NOT NULL,
  instagram VARCHAR(50),
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ativo TINYINT(1) NOT NULL DEFAULT 1
) ENGINE=InnoDB;

-- 3. NIVEIS
CREATE TABLE niveis (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nome VARCHAR(50) NOT NULL UNIQUE
) ENGINE=InnoDB;

-- 4. TIPOS DE CATEGORIA
CREATE TABLE tipos_categoria (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nome VARCHAR(100) NOT NULL UNIQUE
) ENGINE=InnoDB;

-- 5. CATEGORIAS
CREATE TABLE categorias (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_nivel INT NOT NULL,
  descricao VARCHAR(100) NOT NULL,
  id_tipo_categoria INT NOT NULL,
  FOREIGN KEY (id_nivel) REFERENCES niveis(id),
  FOREIGN KEY (id_tipo_categoria) REFERENCES tipos_categoria(id)
) ENGINE=InnoDB;

-- 6. ESPORTES
CREATE TABLE esportes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nome ENUM('Beach Tenis','Tenis de Mesa','Tenis','Pickleball') NOT NULL UNIQUE
) ENGINE=InnoDB;

-- 7. SCOUTS
CREATE TABLE scouts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  vitorias INT NOT NULL DEFAULT 0,
  derrotas INT NOT NULL DEFAULT 0,
  pontos INT NOT NULL DEFAULT 0,
  titulos INT NOT NULL DEFAULT 0
) ENGINE=InnoDB;

-- 8. JOGADORES
CREATE TABLE jogadores (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_usuario INT NOT NULL,
  id_scout INT NOT NULL UNIQUE,
  nome VARCHAR(100) NOT NULL,
  cpf VARCHAR(14) NOT NULL UNIQUE,
  datanascimento DATE NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  telefone VARCHAR(20) NOT NULL,
  whatsapp VARCHAR(20),
  instagram VARCHAR(50),
  sexo ENUM('M','F') NOT NULL,
  equipamento VARCHAR(255),
  tipo ENUM('destro','canhoto') NOT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ativo TINYINT(1) NOT NULL DEFAULT 1,
  FOREIGN KEY (id_usuario) REFERENCES usuarios(id),
  FOREIGN KEY (id_scout) REFERENCES scouts(id)
) ENGINE=InnoDB;

-- 9. JOGADORES_ESPORTES (N:N)
CREATE TABLE jogadores_esportes (
  id_jogador INT NOT NULL,
  id_esporte INT NOT NULL,
  PRIMARY KEY (id_jogador, id_esporte),
  FOREIGN KEY (id_jogador) REFERENCES jogadores(id) ON DELETE CASCADE,
  FOREIGN KEY (id_esporte) REFERENCES esportes(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 10. CLUBES
CREATE TABLE clubes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_jogador_responsavel INT,
  nome VARCHAR(100) NOT NULL,
  telefone VARCHAR(20) NOT NULL,
  whatsapp VARCHAR(20),
  instagram VARCHAR(50),
  estado VARCHAR(50) NOT NULL,
  cidade VARCHAR(100) NOT NULL,
  pais VARCHAR(50) NOT NULL DEFAULT 'Brasil',
  quantidade INT NOT NULL DEFAULT 0,
  ativo TINYINT(1) NOT NULL DEFAULT 1,
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id_jogador_responsavel) REFERENCES jogadores(id)
) ENGINE=InnoDB;

-- 11. CLUBES_USUARIOS (N:N)
CREATE TABLE clubes_usuarios (
  id_clube INT NOT NULL,
  id_jogador INT NOT NULL,
  data_adesao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id_clube, id_jogador),
  FOREIGN KEY (id_clube) REFERENCES clubes(id) ON DELETE CASCADE,
  FOREIGN KEY (id_jogador) REFERENCES jogadores(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 12. TORNEIOS
CREATE TABLE torneios (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_esporte INT NOT NULL,
  nome VARCHAR(100) NOT NULL,
  descricao TEXT,
  quantidade_quadras INT NOT NULL DEFAULT 1,
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id_esporte) REFERENCES esportes(id)
) ENGINE=InnoDB;

-- 13. DUPLAS
CREATE TABLE duplas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_jogador_a INT NOT NULL,
  id_jogador_b INT NOT NULL,
  nome_dupla VARCHAR(100),
  FOREIGN KEY (id_jogador_a) REFERENCES jogadores(id) ON DELETE CASCADE,
  FOREIGN KEY (id_jogador_b) REFERENCES jogadores(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 14. PARTICIPANTES
CREATE TABLE participantes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_torneio INT NOT NULL,
  id_jogador INT,
  id_categoria INT NOT NULL,
  id_dupla INT,
  tipo_modalidade ENUM('simples','duplas') NOT NULL DEFAULT 'simples',
  FOREIGN KEY (id_torneio) REFERENCES torneios(id) ON DELETE CASCADE,
  FOREIGN KEY (id_jogador) REFERENCES jogadores(id) ON DELETE CASCADE,
  FOREIGN KEY (id_categoria) REFERENCES categorias(id),
  FOREIGN KEY (id_dupla) REFERENCES duplas(id)
) ENGINE=InnoDB;

-- 15. GRUPOS
CREATE TABLE grupos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_categoria INT NOT NULL,
  nome VARCHAR(50) NOT NULL,
  FOREIGN KEY (id_categoria) REFERENCES categorias(id)
) ENGINE=InnoDB;

-- 16. GRUPO_PARTICIPANTES (N:N)
CREATE TABLE grupo_participantes (
  id_grupo INT NOT NULL,
  id_participante INT NOT NULL,
  PRIMARY KEY (id_grupo, id_participante),
  FOREIGN KEY (id_grupo) REFERENCES grupos(id) ON DELETE CASCADE,
  FOREIGN KEY (id_participante) REFERENCES participantes(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 17. RODADAS
CREATE TABLE rodadas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nome VARCHAR(50) NOT NULL
) ENGINE=InnoDB;

-- 18. JOGOS
CREATE TABLE jogos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_torneio INT NOT NULL,
  id_grupo INT NOT NULL,
  id_rodada INT NOT NULL,
  id_participante1 INT,
  id_participante2 INT,
  id_dupla1 INT,
  id_dupla2 INT,
  tipo_modalidade ENUM('simples','duplas') NOT NULL DEFAULT 'simples',
  data_hora DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  localizacao VARCHAR(100),
  situacao ENUM('aguardando','em andamento','encerrado') NOT NULL DEFAULT 'aguardando',
  FOREIGN KEY (id_torneio) REFERENCES torneios(id) ON DELETE CASCADE,
  FOREIGN KEY (id_grupo) REFERENCES grupos(id),
  FOREIGN KEY (id_rodada) REFERENCES rodadas(id),
  FOREIGN KEY (id_participante1) REFERENCES participantes(id),
  FOREIGN KEY (id_participante2) REFERENCES participantes(id),
  FOREIGN KEY (id_dupla1) REFERENCES duplas(id),
  FOREIGN KEY (id_dupla2) REFERENCES duplas(id)
) ENGINE=InnoDB;

-- 19. PLACARES
CREATE TABLE placares (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_jogo INT,
  resultado VARCHAR(100),
  id_jogador_vencedor INT,
  id_jogador_perdedor INT,
  FOREIGN KEY (id_jogo) REFERENCES jogos(id) ON DELETE SET NULL,
  FOREIGN KEY (id_jogador_vencedor) REFERENCES jogadores(id) ON DELETE SET NULL,
  FOREIGN KEY (id_jogador_perdedor) REFERENCES jogadores(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- 20. ÍNDICES ÚTEIS
CREATE INDEX idx_jogadores_nome ON jogadores(nome);
CREATE INDEX idx_torneios_nome ON torneios(nome);
CREATE INDEX idx_jogos_data_hora ON jogos(data_hora);
