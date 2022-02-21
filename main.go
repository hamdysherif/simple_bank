package main

import (
	"database/sql"
	"log"

	"github.com/hamdysherif/simplebank/api"
	db "github.com/hamdysherif/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cann't connect to db", err)
		return
	}

	server := api.NewServer(db.NewStore(conn))

	server.Start("0.0.0.0:3009")
}
