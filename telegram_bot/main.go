package main

import (
	"os"

	"github.com/gorilla/mux"
	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pb "github.com/matvoy/chat_server/telegram_bot/proto/bot_message"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel  string
	TgWebhook string
	CertPath  string
	KeyPath   string
	AppPort   int
}

var (
	client  pbstorage.StorageService
	logger  *zerolog.Logger
	cfg     *Config
	service micro.Service
	tgBot   ChatServer
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
			&cli.StringFlag{
				Name:    "tg_webhook_address",
				EnvVars: []string{"TG_WEBHOOK_ADDRESS"},
				Usage:   "Telegram webhook address",
			},
			&cli.IntFlag{
				Name:    "app_port",
				EnvVars: []string{"APP_PORT"},
				Usage:   "Local webhook port",
			},
			&cli.StringFlag{
				Name:    "cert_path",
				EnvVars: []string{"CERT_PATH"},
				Usage:   "SSl certificate",
			},
			&cli.StringFlag{
				Name:    "key_path",
				EnvVars: []string{"KEY_PATH"},
				Usage:   "SSl key",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) error {
			cfg.LogLevel = c.String("log_level")
			cfg.TgWebhook = c.String("tg_webhook_address")
			cfg.CertPath = c.String("cert_path")
			cfg.KeyPath = c.String("key_path")
			cfg.AppPort = c.Int("app_port")
			// cfg.ConversationTimeout = c.Uint64("conversation_timeout")

			client = pbstorage.NewStorageService("webitel.chat.service.storage", service.Client())
			var err error
			logger, err = NewLogger(cfg.LogLevel)
			if err != nil {
				return err
			}
			return configureTelegram()
		}),
		micro.AfterStart(
			func() error {
				return tgBot.StartWebhookServer()
			},
		),
	)

	if err := service.Run(); err != nil {
		logger.Fatal().
			Str("app", "failed to run service").
			Msg(err.Error())
	}
}

func configureTelegram() error {
	r := mux.NewRouter()

	tgBot = NewTelegramBot(
		logger,
		client,
		r,
	)

	if err := pb.RegisterTelegramBotServiceHandler(service.Server(), tgBot); err != nil {
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
