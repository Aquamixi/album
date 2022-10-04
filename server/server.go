package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"main.go/banco"
)

type album struct {
	ID     uint32  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

func BuscaAlbuns(write http.ResponseWriter, request *http.Request) {
	db, err := banco.Conectar()
	if err != nil {
		write.Write([]byte("erro ao conectar"))
		return
	}

	defer db.Close()

	linhas, err := db.Query("select * from album")
	if err != nil {
		write.Write([]byte("erro ao pesquisar"))
		return
	}
	defer linhas.Close()

	var albums []album

	for linhas.Next() {
		var album album

		if err := linhas.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			write.Write([]byte("erro ao escanear"))
			return
		}

		albums = append(albums, album)
	}

	write.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(write).Encode(albums); err != nil {
		write.Write([]byte("erro no encode"))
		return
	}
}

func BuscaAlbun(write http.ResponseWriter, request *http.Request) {
	paramentro := mux.Vars(request)

	ID, err := strconv.ParseUint(paramentro["id"], 10, 32)
	if err != nil {
		write.Write([]byte("erro ao converter id"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		write.Write([]byte("erro ao conectar"))
		return
	}

	linha, err := db.Query("select * from album where id = ?", ID)
	if err != nil {
		write.Write([]byte("erro ao pesquisar"))
		return
	}
	defer db.Close()

	var album album
	if linha.Next() {
		if err := linha.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			write.Write([]byte("erro ao escanear"))
			return
		}
	}

	write.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(write).Encode(album); err != nil {
		write.Write([]byte("erro no encode"))
		return
	}
}

func CriaAlbum(write http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		write.Write([]byte("erro readAll"))
		return
	}

	var album album

	if err := json.Unmarshal(body, &album); err != nil {
		write.Write([]byte("erro Unmarshal"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		write.Write([]byte("erro ao conectar"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("insert into album (title, artist, price) values (?, ?, ?)")
	if err != nil {
		write.Write([]byte("erro ao inserir"))
		return
	}
	defer statement.Close()

	insercao, err := statement.Exec(album.Title, album.Artist, album.Price)
	if err != nil {
		write.Write([]byte("erro ao executar"))
		return
	}

	id_inserido, err := insercao.LastInsertId()
	if err != nil {
		write.Write([]byte("erro ao buscar ultimo id"))
		return
	}

	write.Write([]byte(fmt.Sprintf("ultimo id inserido: %d", id_inserido)))
}

func AtualizaAlbun(write http.ResponseWriter, request *http.Request) {
	parametro := mux.Vars(request)

	ID, err := strconv.ParseUint(parametro["id"], 10, 32)
	if err != nil {
		write.Write([]byte("erro ao pegar o parametro"))
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		write.Write([]byte("erro readAll"))
		return
	}

	var album album

	if err = json.Unmarshal(body, &album); err != nil {
		write.Write([]byte("erro Unmarshal"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		write.Write([]byte("erro ao conectar no banco"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update album set title = ?, artist = ?, price = ? where id = ?")
	if err != nil {
		write.Write([]byte("erro ao atualizar"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(album.Title, album.Artist, album.Price, ID); err != nil {
		write.Write([]byte("erro ao executar"))
		return
	}

	write.WriteHeader(http.StatusNoContent)
}

func DeletaAlbum(write http.ResponseWriter, request *http.Request) {
	parametro := mux.Vars(request)

	ID, err := strconv.ParseUint(parametro["id"], 10, 32)
	if err != nil {
		write.Write([]byte("erro ao pegar o parametro"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		write.Write([]byte("erro ao conectar"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("delete from album where id = ?")
	if err != nil {
		write.Write([]byte("erro ao deletar"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		write.Write([]byte("erro ao executar"))
		return
	}

	write.WriteHeader(http.StatusNoContent)
}
