package main

import (
	"os"

	pb "github.com/matvoy/chat_server/flow_adapter/proto/adapter"
	pbtelegram "github.com/matvoy/chat_server/telegram_bot/proto/bot_message"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel string
}

var (
	client pbtelegram.TelegramBotService
	logger *zerolog.Logger
	cfg    *Config
)

func main() {
	cfg = &Config{}
	service := micro.NewService(
		micro.Name("webitel.chat.service.flowadapter"),
		micro.Version("latest"),
		micro.Flags(
			&cli.StringFlag{
				Name:    "log_level",
				EnvVars: []string{"LOG_LEVEL"},
				Value:   "debug",
				Usage:   "Log Level",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) error {
			cfg.LogLevel = c.String("log_level")
			var err error
			logger, err = NewLogger(cfg.LogLevel)
			client = pbtelegram.NewTelegramBotService("webitel.chat.service.telegrambot", service.Client())
			return err
		}),
	)

	serv := NewFlowService(logger, client)

	if err := pb.RegisterAdapterServiceHandler(service.Server(), serv); err != nil {
		logger.Fatal().
			Str("app", "failed to register service").
			Msg(err.Error())
		return
	}

	if err := service.Run(); err != nil {
		logger.Fatal().
			Str("app", "failed to run service").
			Msg(err.Error())
	}
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
