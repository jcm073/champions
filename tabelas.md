# Tabelas utilizadas pelo sistema.

## Conveções utilizadas

* **Tabelas**: nome sempre em minuscula e no plural
* **Colunas**: nome sempre em minuscula e no singular caso necessário utilize _ (sublinhado/underscore)
* **Chaves Estrangeiras**: utilizar id_xxx exemplo: id_usuarios, id_categorias, etc 


## Tabelas

### usuarios
id

nome

cpf

datanascimento

email

telefone

instagram

datacriação

ativo (true ou false)

### nivel
id

user_id

iniciante

intermediario

intermediarioplus

avançado

profissional

descrição

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
id_usuarios
beachtenis
tenisdemesa
tenis
pickeball

### clubes
id
id_usuarios
nome
estado
cidade
pais
quantidade quadras/mesas)
ativo (true ou false)
datacriação


### torneios
id
nome
descrição
