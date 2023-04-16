package main

import (
	"database/sql"
	"github.com/CrunchyBlue/Golang-Bank/api"
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/CrunchyBlue/Golang-Bank/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
