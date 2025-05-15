package view

import (
	"encoding/json"
	"net/http"
	"champions/presenter"
)

type JogadoresHandler struct {
	Presenter *presenter.JogadoresPresenter
}

func (h *JogadoresHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	jogadores, err := h.Presenter.ListarJogadores()
	if err != nil {
		http.Error(w, "Erro ao buscar jogadores", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jogadores)
}