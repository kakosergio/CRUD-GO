package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// CriarUsuario insere um usuário no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	// cria uma variável pra receber o usuario convertido da request
	var usuario usuario
	// cria uma variavel pra guardar a ID que receberá ao cadastrar o usuario
	var id uint32

	// converte a request body para leitura
	corpoRequisicao, err := io.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição!"))
		return
	}

	// tenta converter o body da requisição em JSON
	if err = json.Unmarshal(corpoRequisicao, &usuario); err != nil {
		w.Write([]byte("Erro ao converter o usuário para struct"))
		return
	}

	// conecta ao banco de dados
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	// já difere o close do banco de dados pra não esquecer
	defer db.Close()

	// cria um statement, um preparo da query do banco de dados só com coringas para incluir os valores depois
	statement, err := db.Prepare("INSERT INTO usuarios (nome, email) values ($1, $2) RETURNING id")
	if err != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	// tem que fechar o statement tbm, então difere pro final
	defer statement.Close()

	// esse eu tive que descobrir na internet kkkk. Usa o statement como base para fazer a query definitiva, passando
	// as informações necessárias para o insert na query, e scan para pegar o id que retorna, pq foi pedido na query
	err = statement.QueryRow(usuario.Nome, usuario.Email).Scan(&id)
	if err != nil {
		w.Write([]byte("Erro ao executar o statement!"))
		return
	}

	// STATUS CODES
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário inserido com sucesso! ID: %d", id)))
}

// BuscarUsuarios traz todos os usuários salvos no banco de dados
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, err := banco.Conectar()

	if err != nil {
		w.Write([]byte("Erro ao conectar ao banco de dados!"))
		return
	}
	defer db.Close()

	linhas, err := db.Query("SELECT * FROM usuarios")

	if err != nil {
		w.Write([]byte("Erro ao buscar os usuários!"))
		return
	}
	defer linhas.Close()

	var usuarios []usuario

	for linhas.Next() {
		var usuario usuario

		if err := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("Erro ao scanear usuário!"))
			return
		}

		usuarios = append(usuarios, usuario)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		w.Write([]byte("Erro ao converter os usuários para JSON!"))
		return
	}
}

// BuscarUsuario traz um usuário específico salvo no banco de dados
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	id, err := strconv.ParseUint(parametros["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao converter parametro para inteiro"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	linha, err := db.Query("SELECT * FROM usuarios WHERE id = $1", id)

	if err != nil {
		w.Write([]byte("Erro ao buscar usuário!"))
		return
	}
	var usuario usuario

	if linha.Next() {
		if err := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("Erro ao scanear as informações do usuário"))
			return
		}
	}
	defer linha.Close()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usuario); err != nil {
		w.Write([]byte("Erro ao converter usuario para JSON"))
		return
	}
}

// AtualizarUsuario atualiza um usuário existente
func AtualizarUsuario(w http.ResponseWriter, r *http.Request){
	parametros := mux.Vars(r)

	id, err := strconv.ParseUint(parametros["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro!"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Erro ao ler o corpo da requisição!"))
		return
	}

	var usuario usuario
	if err := json.Unmarshal(body, &usuario); err != nil {
		w.Write([]byte("Erro ao converter o usuário para struct!"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE usuarios SET nome = $1, email = $2 WHERE id = $3")
	if err != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(usuario.Nome, usuario.Email, id); err != nil {
		w.Write([]byte("Erro ao atualizar o usuário (erro na execução do statement)!"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeletarUsuario apaga um usuário existente
func DeletarUsuario(w http.ResponseWriter, r *http.Request){
	parametros := mux.Vars(r)

	id, err := strconv.ParseUint(parametros["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Erro ao converter parâmetro para int!"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM usuarios WHERE id = $1")
	if err != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(id); err != nil {
		w.Write([]byte("Erro ao deletar o usuário!"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}