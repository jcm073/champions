package model

import (
	"database/sql"
	"time"
)

type Jogador struct {
	ID           int       `json:"id"`
	IdEsportes   int       `json:"id_esportes"`
	IdScouts     int       `json:"id_scouts"`
	Nome         string    `json:"nome"`
	Cpf          string    `json:"cpf"`
	DataNascimento time.Time `json:"datanascimento"`
	Email        string    `json:"email"`
	Telefone     string    `json:"telefone"`
	Whatsup      string    `json:"whatsup"`
	Instagram    string    `json:"instagram"`
	Sexo         string    `json:"sexo"`
	Equipamento  string    `json:"equipamento"`
	Tipo         string    `json:"tipo"`
	CriadoEm     time.Time `json:"criadoem"`
}

// Exemplo de função para buscar todos jogadores
func GetAllJogadores(db *sql.DB) ([]Jogador, error) {
	rows, err := db.Query("SELECT id, id_esportes, id_scouts, nome, cpf, datanascimento, email, telefone, whatsup, instagram, sexo, equipamento, tipo, criadoem FROM jogadores")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jogadores []Jogador
	for rows.Next() {
		var j Jogador
		err := rows.Scan(&j.ID, &j.IdEsportes, &j.IdScouts, &j.Nome, &j.Cpf, &j.DataNascimento, &j.Email, &j.Telefone, &j.Whatsup, &j.Instagram, &j.Sexo, &j.Equipamento, &j.Tipo, &j.CriadoEm)
		if err != nil {
			return nil, err
		}
		jogadores = append(jogadores, j)
	}
	return jogadores, nil
}