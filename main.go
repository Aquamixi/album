package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"main.go/server"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/albuns", server.BuscaAlbuns).Methods(http.MethodGet)
	router.HandleFunc("/album", server.CriaAlbum).Methods(http.MethodPost)
	router.HandleFunc("/album/{id}", server.BuscaAlbun).Methods(http.MethodGet)
	router.HandleFunc("/album/{id}", server.AtualizaAlbun).Methods(http.MethodPut)
	router.HandleFunc("/album/{id}", server.DeletaAlbum).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":5000", router))
}
