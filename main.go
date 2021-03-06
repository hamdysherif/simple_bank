package main

import (
	"database/sql"
	"log"

	"github.com/hamdysherif/simplebank/api"
	db "github.com/hamdysherif/simplebank/db/sqlc"
	"github.com/hamdysherif/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, errC := util.LoadConfig(".")
	if errC != nil {
		log.Fatal("cann't load configs, ", errC)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cann't connect to db", err)
		return
	}

	server, err := api.NewServer(db.NewStore(conn), config)
	if err != nil {
		log.Fatal("cann't create the server", err)
		return
	}

	server.Start(config.ServerAddress)
}
