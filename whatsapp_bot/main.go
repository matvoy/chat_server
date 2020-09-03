package main

import (
	"os"

	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pb "github.com/matvoy/chat_server/whatsapp_bot/proto/bot_message"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel string
}

var (
	client  pbstorage.StorageService
	logger  *zerolog.Logger
	cfg     *Config
	service micro.Service
	wtsBot  ChatServer
)

func init() {
	// plugins
	cmd.DefaultRegistries["consul"] = consul.NewRegistry
}

func main() {
	cfg = &Config{}

	service = micro.NewService(
		micro.Name("webitel.chat.service.telegrambot"),
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

			client = pbstorage.NewStorageService("webitel.chat.service.storage", service.Client())
			var err error
			logger, err = NewLogger(cfg.LogLevel)
			if err != nil {
				return err
			}
			return configureViber()
		}),
		micro.AfterStart(
			func() error {
				return wtsBot.StartServer()
			},
		),
	)

	if err := service.Run(); err != nil {
		logger.Fatal().
			Str("app", "failed to run service").
			Msg(err.Error())
	}
}

func configureViber() error {

	wtsBot = NewWhatsappBotServer(
		logger,
		client,
	)

	if err := pb.RegisterWhatsappBotServiceHandler(service.Server(), wtsBot); err != nil {
		logger.Fatal().
			Str("app", "failed to register service").
			Msg(err.Error())
		return err
	}
	return nil
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
