package main

import (
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/puremike/event-mgt-api/internal/db"
	"github.com/puremike/event-mgt-api/internal/env"
	"github.com/puremike/event-mgt-api/internal/storage"
	"go.uber.org/zap"
)

type application struct {
	config *config
	store  *storage.Storage
	logger *zap.SugaredLogger
}

type config struct {
	port     string
	env      string
	dbconfig dbconfig
}

type dbconfig struct {
	db_url                     string
	maxIdleConns, maxOpenConns int
	connMaxIdleTime            time.Duration
}

func main() {

	cfg := &config{port: env.GetEnvString("PORT", "5300"), env: env.GetEnvString("ENV", "development"),
		dbconfig: dbconfig{
			db_url: env.GetEnvString("DB_URL", "postgres://postgres:postgress123@localhost/EventMGTAPI?sslmode=disable"), maxIdleConns: env.GetEnvInt("SET_MAX_IDLE_CONNS", 8), maxOpenConns: env.GetEnvInt("SET_MAX_OPEN_CONNS", 50), connMaxIdleTime: env.GetEnvDuration("SET_CONN_MAX_IDLE_TIME", 20*time.Minute),
		}}

	db, err := db.ConnectPostgresDB(cfg.dbconfig.db_url, cfg.dbconfig.maxIdleConns, cfg.dbconfig.maxOpenConns, cfg.dbconfig.connMaxIdleTime)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	logger.Info("DB connection opened successfully")

	app := application{
		config: cfg,
		store:  storage.NewStorage(db),
		logger: logger,
	}

	mux := app.routes()
	log.Fatal(app.server(mux))
}
