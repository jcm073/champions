// main.go
package main

import (
	"log"
	"net/http"

	"champions/view"
)

func main() {
	router := view.NewRouter()
	log.Println("API rodando em :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}