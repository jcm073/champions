# Tabelas utilizadas pelo sistema.

## Conveções utilizadas

* **Tabelas**: nome sempre em minuscula e no plural
* **Colunas**: nome sempre em minuscula e no singular caso necessário utilize _ (sublinhado/underscore)
* **Chaves Estrangeiras**: utilizar id_xxx exemplo: id_usuarios, id_categorias, etc 


## Tabelas

### usuarios
id
id_perfil
nome
cpf
datanascimento
email
telefone
instagram
criadoem
ativo (true ou false)


### perfis
id
tipo {
    jogador
    usuario
    admim
    gestor_clube
    gestor_torneio
}

### jogadores
id
id_esportes
id_scouts
nome
cpf
datanascimento
email
telefone
whatsup
instagram
sexo
equipamento (marca e modelos da raquete)
tipo (destro ou canhoto)
criadoem

### niveis
id
descrição
iniciante
intermediario
intermediarioplus
avançado
profissional

### categorias
id
id_nivel
descrição
simplesmasculina
simplesfeminina
duplasmasculina
duplasfeminina
duplasmista
40feminina
40masculina
50feminina
50masculina
60feminina
60masculina
veteranosfeminino
veteranosmasculino


### esportes
id
name (beachtenis, tenisdemesa, tenis, pickeball)

### clubes
id
id_jogadores
nome
responsavel
telefone
whatsup
instagram
estado
cidade
pais
quantidade (quadras/mesas)
ativo (true ou false)
criadoem

### clubes_usuarios
id_clubes
id_jogadores
descrição
datainserção
criadoem

### torneios
id
id_esportes
nome
descrição
quantidadequadras
criadoem

### participantes (para saber todos os jogadores q estão inscritos no torneio e as devidas categorias)
id
id_torneios
id_jogadores
id_categorias

### grupos (mínimo 2 e máximo 5 jogadores)
id
id_categoria
name (a, b, c, ou 1, 2, 3)
id_jogador1
id_jogador2
id_jogador3
id_jogador4

### rodadas
id
name (1,2,3, oitava-finais, quartas-finais, semi-finais, final)

### jogos
id
id_torneios
id_grupos
id_rodadas
id_jogador1
id_jogador2
id_placares
datetime
localização (quadra1, quadra2, etc)
situação (aguardando inicio, em andamento, encerrado)

### placares
id
id_jogos
pontuação_jogador1
pontuação_jogador2
resultado (1x2)
vencedor (jogador2)
perdedor (jogador1)

### scouts
id
id_players
id_esporte
id_torneios
vitorias
derrotas
pontos
titulos