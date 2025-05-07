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
  tipo ENUM('jogador', 'usuario', 'admim', 'gestor_clube', 'gestor_torneio') NOT NULL
) ENGINE=InnoDB;

-- Tabela usuarios
CREATE TABLE usuarios (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_perfil INT NOT NULL,
  nome VARCHAR(100) NOT NULL,
  cpf VARCHAR(14) NOT NULL,
  datanascimento DATE NOT NULL,
  email VARCHAR(100) NOT NULL,
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
  name ENUM('beachtenis', 'tenisdemesa', 'tenis', 'pickeball') NOT NULL
) ENGINE=InnoDB;

-- Tabela scouts (precisa ser criada antes de jogadores devido à referência circular)
CREATE TABLE scouts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  vitorias INT NOT NULL DEFAULT 0,
  derrotas INT NOT NULL DEFAULT 0,
  pontos INT NOT NULL DEFAULT 0,
  titulos INT NOT NULL DEFAULT 0
) ENGINE=InnoDB;

-- Tabela jogadores
CREATE TABLE jogadores (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_esportes INT NOT NULL,
  id_scouts INT NOT NULL,
  nome VARCHAR(100) NOT NULL,
  cpf VARCHAR(14) NOT NULL,
  datanascimento DATE NOT NULL,
  email VARCHAR(100) NOT NULL,
  telefone VARCHAR(20) NOT NULL,
  whatsup VARCHAR(20),
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

-- Atualizar a tabela scouts para adicionar as chaves estrangeiras
ALTER TABLE scouts
  ADD COLUMN id_players INT,
  ADD COLUMN id_esporte INT,
  ADD COLUMN id_torneios INT,
  ADD FOREIGN KEY (id_players) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL,
  ADD FOREIGN KEY (id_esporte) REFERENCES esportes(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL;

-- Tabela clubes
CREATE TABLE clubes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_jogadores INT,
  nome VARCHAR(100) NOT NULL,
  responsavel VARCHAR(100) NOT NULL,
  telefone VARCHAR(20) NOT NULL,
  whatsup VARCHAR(20),
  instagram VARCHAR(50),
  estado VARCHAR(50) NOT NULL,
  cidade VARCHAR(100) NOT NULL,
  pais VARCHAR(50) NOT NULL DEFAULT 'Brasil',
  quantidade INT NOT NULL DEFAULT 0,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id_jogadores) REFERENCES jogadores(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
) ENGINE=InnoDB;

-- Tabela clubes_usuarios (tabela de junção)
CREATE TABLE clubes_usuarios (
  id_clubes INT NOT NULL,
  id_jogadores INT NOT NULL,
  descrição VARCHAR(255),
  datainserção DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
  descrição VARCHAR(100) NOT NULL,
  iniciante BOOLEAN NOT NULL DEFAULT FALSE,
  intermediario BOOLEAN NOT NULL DEFAULT FALSE,
  intermediarioplus BOOLEAN NOT NULL DEFAULT FALSE,
  avançado BOOLEAN NOT NULL DEFAULT FALSE,
  profissional BOOLEAN NOT NULL DEFAULT FALSE
) ENGINE=InnoDB;

-- Tabela categorias
CREATE TABLE categorias (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_nivel INT NOT NULL,
  descrição VARCHAR(100) NOT NULL,
  simplesmasculina BOOLEAN NOT NULL DEFAULT FALSE,
  simplesfeminina BOOLEAN NOT NULL DEFAULT FALSE,
  duplasmasculina BOOLEAN NOT NULL DEFAULT FALSE,
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
    ON DELETE RESTRICT
) ENGINE=InnoDB;

-- Tabela torneios
CREATE TABLE torneios (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_esportes INT NOT NULL,
  nome VARCHAR(100) NOT NULL,
  descrição TEXT,
  quantidadequadras INT NOT NULL DEFAULT 1,
  criadoem DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id_esportes) REFERENCES esportes(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
) ENGINE=InnoDB;

-- Completar a referência circular em scouts
ALTER TABLE scouts
  ADD FOREIGN KEY (id_torneios) REFERENCES torneios(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL;

-- Tabela participantes
CREATE TABLE participantes (
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
) ENGINE=InnoDB;

-- Tabela grupos
CREATE TABLE grupos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_categoria INT NOT NULL,
  name VARCHAR(50) NOT NULL,
  id_jogador1 INT NOT NULL,
  id_jogador2 INT NOT NULL,
  id_jogador3 INT,
  id_jogador4 INT,
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
  name VARCHAR(50) NOT NULL
) ENGINE=InnoDB;

-- Tabela placares (precisa ser criada antes de jogos devido à referência circular)
CREATE TABLE placares (
  id INT AUTO_INCREMENT PRIMARY KEY,
  pontuação_jogador1 INT NOT NULL DEFAULT 0,
  pontuação_jogador2 INT NOT NULL DEFAULT 0,
  resultado VARCHAR(100),
  vencedor VARCHAR(100),
  perdedor VARCHAR(100)
) ENGINE=InnoDB;

-- Tabela jogos
CREATE TABLE jogos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_torneios INT NOT NULL,
  id_grupos INT NOT NULL,
  id_rodadas INT NOT NULL,
  id_jogador1 INT NOT NULL,
  id_jogador2 INT NOT NULL,
  id_placares INT NOT NULL,
  datetime DATETIME NOT NULL,
  localização VARCHAR(100),
  situação ENUM('aguardando', 'em andamento', 'encerrado') NOT NULL DEFAULT 'aguardando',
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
  ADD FOREIGN KEY (id_jogos) REFERENCES jogos(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL;