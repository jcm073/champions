package view

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"champions/presenter"
)

func NewRouter() http.Handler {
	// Configure o seu DSN conforme necessário
	db, err := sql.Open("mysql", "usuario:senha@tcp(localhost:3306)/campeonatos")
	if err != nil {
		panic(err)
	}

	jogadoresPresenter := &presenter.JogadoresPresenter{DB: db}
	jogadoresHandler := &JogadoresHandler{Presenter: jogadoresPresenter}

	router := mux.NewRouter()
	router.HandleFunc("/jogadores", jogadoresHandler.GetAll).Methods("GET")
	// Adicione rotas para outras entidades aqui

	return router
}