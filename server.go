package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	Id    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var usuario usuario

	if err = json.Unmarshal(req, &usuario); err != nil {
		w.Write([]byte("Erro ao converter usuário para struct"))
		return
	}

	db, err := database()
	if err != nil {
		w.Write([]byte("Falha ao conectar com o banco de dados"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("insert into usuarios (nome, email) values (?, ?)")
	if err != nil {
		w.Write([]byte("Falha ao criar o statement"))
		return
	}
	defer statement.Close()

	insertion, err := statement.Exec(usuario.Nome, usuario.Email)
	if err != nil {
		w.Write([]byte("Falha ao executar o statement"))
		return
	}

	idInserted, err := insertion.LastInsertId()
	if err != nil {
		w.Write([]byte("Falha ao obter id inserido"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário inserido com sucessl! Id: %d", idInserted)))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database()
	if err != nil {
		w.Write([]byte("Falha ao conectar com o banco de dados"))
		return
	}
	defer db.Close()

	registries, err := db.Query("select * from usuarios")
	if err != nil {
		w.Write([]byte("Falha ao buscar os usuários"))
		return
	}
	defer registries.Close()

	var usuarios []usuario
	for registries.Next() {
		var usuario usuario

		if err := registries.Scan(&usuario.Id, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("Falha ao escanear o usuario"))
			return
		}

		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		w.Write([]byte("Erro ao converter os usuarios para JSON"))
		return
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Falha converter o parametro para inteiro"))
		return
	}

	db, err := database()
	if err != nil {
		w.Write([]byte("Falha ao conectar com o banco de dados"))
		return
	}
	defer db.Close()

	registry, err := db.Query("select * from usuarios where id =  ?", ID)
	if err != nil {
		w.Write([]byte("Falha ao buscar o usuário"))
		return
	}

	var usuario usuario
	if registry.Next() {
		if err := registry.Scan(&usuario.Id, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("Falha ao escanear o usuario"))
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usuario); err != nil {
		w.Write([]byte("Erro ao converter o usuario para JSON"))
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Falha converter o parametro para inteiro"))
		return
	}

	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var usuario usuario
	if err = json.Unmarshal(req, &usuario); err != nil {
		w.Write([]byte("Erro ao converter usuário para struct"))
		return
	}

	db, err := database()
	if err != nil {
		w.Write([]byte("Falha ao conectar com o banco de dados"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update usuarios set nome=?, email=? where id=?")
	if err != nil {
		w.Write([]byte("Falha preparar o statement"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(usuario.Nome, usuario.Email, ID); err != nil {
		w.Write([]byte("Erro ao atualizar o usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Falha converter o parametro para inteiro"))
		return
	}

	db, err := database()
	if err != nil {
		w.Write([]byte("Falha ao conectar com o banco de dados"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("delete from usuarios where id=?")
	if err != nil {
		w.Write([]byte("Falha preparar o statement"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		w.Write([]byte("Erro ao deletar o usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
