package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/mikosco4real/simple_bank/api"
	db "github.com/mikosco4real/simple_bank/db/sqlc"
	"github.com/mikosco4real/simple_bank/util"
)


func main() {
	config, err := util.LoadConfig(".")
	
	if err != nil {
		log.Fatal("Cannot load config values: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start http server on ", config.ServerAddress)
	}
}