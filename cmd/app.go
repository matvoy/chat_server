package main

import (
	"os"

	"github.com/matvoy/chat_server/app/chats"
	"github.com/matvoy/chat_server/app/chats/telegram"

	"github.com/rs/zerolog"
)

type App struct {
	config        *Config
	log           *zerolog.Logger
	telegramBot   chats.ChatServer
	websocketChat chats.ChatServer
}

func NewApp(cfg *Config) (*App, error) {
	logger := newLogger(cfg.Debug)
	a := &App{
		config:      cfg,
		log:         logger,
		telegramBot: telegram.NewTelegramBot(cfg.TelegramBotToken, logger),
	}

	return a, nil
}

func (a *App) Start() {
	a.log.Info().
		Str("app", "start listening").
		Int("port", a.config.AppPort).
		Msg("press Ctrl-C to shutdown")
	if err := a.telegramBot.Start(); err != nil {
		a.log.Fatal().
			Str("app", "failed to start").
			Msg(err.Error())
	}

}

func (a *App) Stop() {
	a.log.Info().
		Msg("app: stopping")
}

func newLogger(isDebug bool) *zerolog.Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(logLevel)
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &l
}
