package banco

import (
	"database/sql"

	_ "github.com/lib/pq" // driver de conexão ao bd postgres
)

// Conectar abre a conexão com o bd
func Conectar() (*sql.DB, error) {
	// String que guarda as informações para conexão no banco de dados
	stringConn := "host=localhost port=5432 user=golang password=arthasmyson dbname=devbook sslmode=disable"
	// abre a conexão com o banco, indicando driver e a string com informações de conexão
	db, err := sql.Open("postgres", stringConn)
	// se tiver erro, retorna o erro
	if err != nil {
		return nil, err
	}
	// Testa a conexão com o banco de dados. Se der erro, retorna o erro
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// Se tudo estiver OK, retorna a conexão com o banco de dados aberta
	return db, nil

}
