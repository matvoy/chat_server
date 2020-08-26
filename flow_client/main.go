package main

import (
	"os"

	pbstorage "github.com/matvoy/chat_server/chat_storage/proto/storage"
	pb "github.com/matvoy/chat_server/flow_client/proto/flow_client"
	pbmanager "github.com/matvoy/chat_server/flow_client/proto/flow_manager"
	pbtelegram "github.com/matvoy/chat_server/telegram_bot/proto/bot_message"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/store/redis/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel string
}

var (
	botClient     pbtelegram.TelegramBotService
	managerClient pbmanager.FlowChatServerService
	storageClient pbstorage.StorageService
	logger        *zerolog.Logger
	cfg           *Config
	redisStore    store.Store
	redisTable    string
	timeout       uint64
)

func init() {
	// plugins
	cmd.DefaultStores["redis"] = redis.NewStore
	cmd.DefaultRegistries["consul"] = consul.NewRegistry
}

func main() {
	cfg = &Config{}
	service := micro.NewService(
		micro.Name("webitel.chat.service.flowclient"),
		micro.Version("latest"),
		micro.Flags(
			&cli.StringFlag{
				Name:    "log_level",
				EnvVars: []string{"LOG_LEVEL"},
				Value:   "debug",
				Usage:   "Log Level",
			},
			&cli.Uint64Flag{
				Name:    "conversation_timeout_sec",
				EnvVars: []string{"CONVERSATION_TIMEOUT_SEC"},
				Usage:   "Conversation timeout. sec",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) error {
			cfg.LogLevel = c.String("log_level")
			redisTable = c.String("store_table")
			timeout = 600 //c.Uint64("conversation_timeout_sec")
			var err error
			logger, err = NewLogger(cfg.LogLevel)
			botClient = pbtelegram.NewTelegramBotService("webitel.chat.service.telegrambot", service.Client())
			managerClient = pbmanager.NewFlowChatServerService("workflow", service.Client())
			storageClient = pbstorage.NewStorageService("webitel.chat.service.storage", service.Client())
			return err
		}),
	)

	service.Options().Store.Init(store.Table(redisTable))

	serv := NewFlowService(logger, botClient, managerClient, service.Options().Store, storageClient)

	if err := pb.RegisterFlowAdapterServiceHandler(service.Server(), serv); err != nil {
		logger.Fatal().
			Str("app", "failed to register service").
			Msg(err.Error())
		return
	}

	if err := pb.RegisterFlowClientServiceHandler(service.Server(), serv); err != nil {
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
