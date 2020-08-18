package main

import (
	"os"

	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pbflow "github.com/matvoy/chat_server/flow_adapter/proto/adapter"
	pb "github.com/matvoy/chat_server/telegram_bot/proto/bot_message"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-plugins/store/redis/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel            string
	TelegramBotToken    string
	ProfileID           uint64
	RedisURL            string
	ConversationTimeout uint64
}

const (
	redisTable = "chat:"
)

var (
	client     pbstorage.StorageService
	flowClient pbflow.AdapterService
	logger     *zerolog.Logger
	cfg        *Config
	service    micro.Service
	redisStore store.Store
)

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
				Name:    "telegram_bot_token",
				EnvVars: []string{"TELEGRAM_BOT_TOKEN"},
				Usage:   "Telegram bot token",
			},
			&cli.Uint64Flag{
				Name:    "profile_id",
				EnvVars: []string{"PROFILE_ID"},
				Usage:   "Profile id",
			},
			&cli.StringFlag{
				Name:    "redis_url",
				EnvVars: []string{"REDIS_URL"},
				Usage:   "Redis URL",
			},
			&cli.Uint64Flag{
				Name:    "conversation_timeout",
				EnvVars: []string{"CONVERSATION_TIMEOUT"},
				Usage:   "Conversation timeout",
			},
		),
	)

	redisStore = redis.NewStore(
		store.Nodes("redis://10.9.8.111:6379"),
		store.Table(redisTable),
	)

	service.Init(
		micro.Action(func(c *cli.Context) error {
			cfg.LogLevel = c.String("log_level")
			cfg.TelegramBotToken = c.String("telegram_bot_token")
			cfg.ProfileID = c.Uint64("profile_id")
			cfg.RedisURL = c.String("redis_url")
			cfg.ConversationTimeout = c.Uint64("conversation_timeout")

			client = pbstorage.NewStorageService("webitel.chat.service.storage", service.Client())
			flowClient = pbflow.NewAdapterService("webitel.chat.service.flowadapter", service.Client())
			var err error
			logger, err = NewLogger(cfg.LogLevel)
			return err
		}),
		micro.AfterStart(
			startTelegram,
		),
		micro.Store(redisStore),
	)

	if err := service.Run(); err != nil {
		logger.Fatal().
			Str("app", "failed to run service").
			Msg(err.Error())
	}
}

func startTelegram() error {
	tgBot := NewTelegramBot(
		cfg.TelegramBotToken,
		cfg.ProfileID,
		logger,
		client,
		flowClient,
		redisStore,
		cfg.ConversationTimeout,
	)
	if err := pb.RegisterTelegramBotServiceHandler(service.Server(), tgBot); err != nil {
		logger.Fatal().
			Str("app", "failed to register service").
			Msg(err.Error())
		return err
	}
	return tgBot.Start()
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
