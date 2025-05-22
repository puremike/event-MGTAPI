package main

import (
	"expvar"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/puremike/event-mgt-api/docs"
	"github.com/puremike/event-mgt-api/internal/auth"
	"github.com/puremike/event-mgt-api/internal/db"
	"github.com/puremike/event-mgt-api/internal/env"
	"github.com/puremike/event-mgt-api/internal/storage"
	"github.com/puremike/event-mgt-api/internal/storage/cache"
	"go.uber.org/zap"
)

type application struct {
	config           *config
	store            *storage.Storage
	logger           *zap.SugaredLogger
	jWTAuthenticator *auth.JWTAuthenticator
	cacheStorage     *cache.CacheStorage
}

type config struct {
	port              string
	env               string
	dbconfig          dbconfig
	authConfig        authConfig
	redisClientConfig redisClientConfig
}

type redisClientConfig struct {
	addr, pw string
	db       int
	enabled  bool
}

type authConfig struct {
	secretKey, iss, aud string
	tokenExp            time.Duration
	username, password  string
}

type dbconfig struct {
	db_url                     string
	maxIdleConns, maxOpenConns int
	connMaxIdleTime            time.Duration
}

//	@title			Event Management API
//	@version		1.0
//	@description	This is an API for event management

//	@contact.name	Puremike
//	@contact.url	http://github.com/puremike
//	@contact.email	digitalmarketfy@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Use a valid JWT token. Format: Bearer <token>

func main() {

	docs.SwaggerInfo.BasePath = "/api/v1"

	cfg := &config{port: env.GetEnvString("PORT", "5300"), env: env.GetEnvString("ENV", "development"),
		dbconfig: dbconfig{
			db_url: env.GetEnvString("DB_URL", "postgres://postgres:postgress123@localhost/EventMGTAPI?sslmode=disable"), maxIdleConns: env.GetEnvInt("SET_MAX_IDLE_CONNS", 8), maxOpenConns: env.GetEnvInt("SET_MAX_OPEN_CONNS", 50), connMaxIdleTime: env.GetEnvDuration("SET_CONN_MAX_IDLE_TIME", 20*time.Minute),
		}, authConfig: authConfig{
			secretKey: env.GetEnvString("JWT_SECRET_KEY", "HKSHD7923B799B08409023N988"), iss: env.GetEnvString("JWT_ISS", "event-mgt-api"), aud: env.GetEnvString("JWT_AUD", "event-mgt-api"), tokenExp: env.GetEnvDuration("JWT_TOKEN_EXP", 72*time.Hour), username: env.GetEnvString("BASIC_AUTH_USERNAME", "admin"), password: env.GetEnvString("BASIC_AUTH_PASSWORD", "password")},
		redisClientConfig: redisClientConfig{
			addr:    env.GetEnvString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetEnvString("REDIS_PW", ""),
			db:      env.GetEnvInt("REDIS_DB", 0),
			enabled: env.GetEnvBool("REDIS_ENABLED", false)},
	}

	db, err := db.ConnectPostgresDB(cfg.dbconfig.db_url, cfg.dbconfig.maxIdleConns, cfg.dbconfig.maxOpenConns, cfg.dbconfig.connMaxIdleTime)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	logger.Info("DB connection opened successfully")

	var rdb *redis.Client
	if cfg.redisClientConfig.enabled {

		rdb = cache.NewRedisClient(cfg.redisClientConfig.addr, cfg.redisClientConfig.pw, cfg.redisClientConfig.db)

		logger.Info("Redis connection opened successfully")
	}

	app := application{
		config:           cfg,
		store:            storage.NewStorage(db),
		logger:           logger,
		jWTAuthenticator: auth.NewJWTAuthenticator(cfg.authConfig.secretKey, cfg.authConfig.iss, cfg.authConfig.aud),
		cacheStorage:     cache.NewCacheStorage(rdb),
	}

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	mux := app.routes()
	log.Fatal(app.server(mux))
}
