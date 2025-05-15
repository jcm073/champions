package presenter

import (
	"database/sql"
	"champions/model"
)

type JogadoresPresenter struct {
	DB *sql.DB
}

func (p *JogadoresPresenter) ListarJogadores() ([]model.Jogador, error) {
	return model.GetAllJogadores(p.DB)
}