package env

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT, DB_URL                           string
	SET_MAX_IDLE_CONNS, SET_MAX_OPEN_CONNS int
	SET_CONN_MAX_IDLE_TIME                 time.Duration
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":5300"
	}

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("DB_URL not set")
	}

	maxIdleConns, err := strconv.Atoi(os.Getenv("SET_MAX_IDLE_CONNS"))
	if err != nil {
		log.Println("Invalid SET_MAX_IDLE_CONNS value, defaulting to 10")
		maxIdleConns = 10
	}

	maxOpenConns, err := strconv.Atoi(os.Getenv("SET_MAX_OPEN_CONNS"))
	if err != nil {
		log.Println("Invalid SET_MAX_OPEN_CONNS value, defaulting to 100")
		maxOpenConns = 100
	}

	connMaxIdleTime, err := time.ParseDuration(os.Getenv("SET_CONN_MAX_IDLE_TIME"))
	if err != nil {
		log.Println("Invalid SET_CONNS_MAX_IDLE_TIME value, defaulting to 40m")
		connMaxIdleTime = 40 * time.Minute
	}

	return &Config{
		PORT:                   ":" + port,
		DB_URL:                 db_url,
		SET_MAX_IDLE_CONNS:     maxIdleConns,
		SET_MAX_OPEN_CONNS:     maxOpenConns,
		SET_CONN_MAX_IDLE_TIME: connMaxIdleTime,
	}
}

// func mustAtoi(s string) int {
// 	i, err := strconv.Atoi(s)
// 	if err != nil {
// 		panic("Invalid integer: " + s)
// 	}
// 	return i
// }
