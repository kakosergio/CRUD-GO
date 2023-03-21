package main

import (
	"crud/servidor"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// CRUD - CREATE, READ, UPDATE, DELETE

	// CREATE - POST
	// READ - GET
	// UPDATE - PUT
	// DELETE - DELETE

	// Cria um router com o pacote Mux, para configurações de rota e métodos HTTP
	router := mux.NewRouter()

	// rota usuários, função CriarUsuario, método POST
	router.HandleFunc("/usuarios", servidor.CriarUsuario).Methods(http.MethodPost)
	// rota usuarios, função BuscarUsuarios, método GET
	router.HandleFunc("/usuarios", servidor.BuscarUsuarios).Methods(http.MethodGet)
	// rota usuarios, parâmetro id, função BuscarUsuario, método GET
	router.HandleFunc("/usuarios/{id}", servidor.BuscarUsuario).Methods(http.MethodGet)
	// rota usuarios, parâmetro id, função AtualizarUsuario, método PUT
	router.HandleFunc("/usuarios/{id}", servidor.AtualizarUsuario).Methods(http.MethodPut)
	// rota usuarios, parâmetro id, função DeletarUsuario, método DELETE
	router.HandleFunc("/usuarios/{id}", servidor.DeletarUsuario).Methods(http.MethodDelete)

	// printa no terminal que o servidor está up
	fmt.Println("Escutando na porta 5000")
	// inicia o servidor
	log.Fatal(http.ListenAndServe(":5000", router))
}
