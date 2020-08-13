package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/matvoy/chat_server/chat_api/handlers"
	"github.com/matvoy/chat_server/chat_api/repo/pg"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type App struct {
	config  *Config
	log     *zerolog.Logger
	db      *sql.DB
	handler http.Handler
}

func NewApp(cfg *Config) (*App, error) {
	logger, err := NewLogger(cfg.LogLevel)
	if err != nil {
		logger.Fatal().
			Str("app", "failed to parse log level").
			Msg(err.Error())
		return nil, err
	}
	db, err := sql.Open("postgres", DbSource(cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode))
	if err != nil {
		logger.Fatal().
			Str("app", "failed to connect db").
			Msg(err.Error())
		return nil, err
	}
	repo := pg.NewPgRepository(db, logger)
	handlers := handlers.NewApiHandlers(repo, logger)
	a := &App{
		config: cfg,
		log:    logger,
		db:     db,
	}
	a.registerRoutes(handlers)

	return a, nil
}

func (a *App) Start() {
	a.log.Info().
		Str("app", "start listening").
		Int("port", a.config.ApiHttpPort).
		Msg("press Ctrl-C to shutdown")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", a.config.ApiHttpPort), a.handler); err != nil {
		a.log.Fatal().
			Str("app", "failed to start").
			Msg(err.Error())
	}

}

func (a *App) Stop() {
	a.log.Info().
		Msg("app: stopping")
	if a.db == nil {
		return
	}
	if err := a.db.Close(); err != nil {
		a.log.Error().
			Msg("app: failed to close db connection")
	}
}

func (a *App) registerRoutes(handlers handlers.ApiHandlers) {
	r := mux.NewRouter()
	r.HandleFunc("/profiles", handlers.GetProfiles).Methods("GET")
	r.HandleFunc("/conversations", handlers.GetConversations).Methods("GET")
	r.HandleFunc("/messages", handlers.GetMessages).Methods("GET")
	r.HandleFunc("/clients", handlers.GetClients).Methods("GET")
	r.HandleFunc("/user_conversations", handlers.GetUserConversations).Methods("GET")
	r.HandleFunc("/attachments", handlers.GetAttachments).Methods("GET")
	a.handler = r
}

func DbSource(host, user, dbName, password, sslMode string) string {
	dbinfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s", host, user, dbName, sslMode)
	if password != "" {
		dbinfo += fmt.Sprintf(" password=%s", password)
	}
	return dbinfo
}

func NewLogger(logLevel string) (*zerolog.Logger, error) {
	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	l := zerolog.New(os.Stdout)
	l = l.Level(lvl)

	return &l, nil
}
