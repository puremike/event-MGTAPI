package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/puremike/event-mgt-api/internal/db"
	"github.com/puremike/event-mgt-api/internal/env"
)

var getEnv *env.Config

func main() {
	getEnv = env.Load()

	db, err := db.ConnectPostgresDB(getEnv.PORT, getEnv.SET_MAX_IDLE_CONNS, getEnv.SET_MAX_OPEN_CONNS, getEnv.SET_CONN_MAX_IDLE_TIME)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
