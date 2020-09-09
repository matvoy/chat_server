package main

import (
	"os"

	pbbot "github.com/matvoy/chat_server/api/proto/bot"
	pbchat "github.com/matvoy/chat_server/api/proto/chat"
	pb "github.com/matvoy/chat_server/api/proto/flow_client"
	pbmanager "github.com/matvoy/chat_server/api/proto/flow_manager"
	cache "github.com/matvoy/chat_server/internal/chat_cache"

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
	botClient     pbbot.BotService
	managerClient pbmanager.FlowChatServerService
	chatClient    pbchat.ChatService
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
		micro.Name("webitel.chat.flowclient"),
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
			botClient = pbbot.NewBotService("webitel.chat.bot", service.Client())
			managerClient = pbmanager.NewFlowChatServerService("workflow", service.Client())
			chatClient = pbchat.NewChatService("webitel.chat.server", service.Client())
			return err
		}),
	)

	service.Options().Store.Init(store.Table(redisTable))

	cache := cache.NewChatCache(service.Options().Store)
	serv := NewFlowService(logger, botClient, managerClient, chatClient, cache)

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
