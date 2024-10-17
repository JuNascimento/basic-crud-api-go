package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/usuarios", createUser).Methods(http.MethodPost)
	router.HandleFunc("/usuarios", getUsers).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", getUserById).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", updateUser).Methods(http.MethodPut)
	router.HandleFunc("/usuarios/{id}", deleteUser).Methods(http.MethodDelete)

	fmt.Println("Servindo a aplicação no endereço http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
