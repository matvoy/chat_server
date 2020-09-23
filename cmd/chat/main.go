package main

import (
	"database/sql"
	"fmt"
	"os"

	pbbot "github.com/matvoy/chat_server/api/proto/bot"
	pb "github.com/matvoy/chat_server/api/proto/chat"
	pbflow "github.com/matvoy/chat_server/api/proto/flow_client"
	cache "github.com/matvoy/chat_server/internal/chat_cache"
	"github.com/matvoy/chat_server/internal/repo/pg"

	_ "github.com/lib/pq"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/store/redis/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel   string
	DBHost     string
	DBUser     string
	DBName     string
	DBSSLMode  string
	DBPassword string
}

var (
	logger     *zerolog.Logger
	cfg        *Config
	service    micro.Service
	redisStore store.Store
	// rabbitBroker broker.Broker
	redisTable string
	flowClient pbflow.FlowAdapterService
	botClient  pbbot.BotService
)

func init() {
	// plugins
	cmd.DefaultBrokers["rabbitmq"] = rabbitmq.NewBroker
	cmd.DefaultStores["redis"] = redis.NewStore
	cmd.DefaultRegistries["consul"] = consul.NewRegistry
}

func main() {
	cfg = &Config{}
	service = micro.NewService(
		micro.Name("webitel.chat.server"),
		micro.Version("latest"),
		micro.Flags(
			&cli.StringFlag{
				Name:    "log_level",
				EnvVars: []string{"LOG_LEVEL"},
				Value:   "debug",
				Usage:   "Log Level",
			},
			&cli.StringFlag{
				Name:    "db_host",
				EnvVars: []string{"DB_HOST"},
				Usage:   "DB Host",
			},
			&cli.StringFlag{
				Name:    "db_user",
				EnvVars: []string{"DB_USER"},
				Usage:   "DB User",
			},
			&cli.StringFlag{
				Name:    "db_name",
				EnvVars: []string{"DB_NAME"},
				Usage:   "DB Name",
			},
			&cli.StringFlag{
				Name:    "db_sslmode",
				EnvVars: []string{"DB_SSLMODE"},
				Value:   "disable",
				Usage:   "DB SSL Mode",
			},
			&cli.StringFlag{
				Name:    "db_password",
				EnvVars: []string{"DB_PASSWORD"},
				Usage:   "DB Password",
			},
		),
	)
	service.Init(
		micro.Action(func(c *cli.Context) error {
			cfg.LogLevel = c.String("log_level")
			cfg.DBHost = c.String("db_host")
			cfg.DBUser = c.String("db_user")
			cfg.DBName = c.String("db_name")
			cfg.DBSSLMode = c.String("db_sslmode")
			cfg.DBPassword = c.String("db_password")
			redisTable = c.String("store_table")
			var err error
			logger, err = NewLogger(cfg.LogLevel)
			if err != nil {
				logger.Fatal().
					Str("app", "failed to parse log level").
					Msg(err.Error())
				return err
			}
			flowClient = pbflow.NewFlowAdapterService("webitel.chat.flowclient", service.Client())
			botClient = pbbot.NewBotService("webitel.chat.bot", service.Client())
			return nil
		}),
		micro.Broker(
			rabbitmq.NewBroker(
				rabbitmq.ExchangeName("chat"),
				rabbitmq.DurableExchange(),
			),
		),
	)

	service.Options().Store.Init(store.Table(redisTable))

	if err := service.Options().Broker.Init(); err != nil {
		logger.Fatal().
			Str("app", "failed to init broker").
			Msg(err.Error())
		return
	}
	if err := service.Options().Broker.Connect(); err != nil {
		logger.Fatal().
			Str("app", "failed to connect broker").
			Msg(err.Error())
		return
	}

	db, err := sql.Open("postgres", DbSource(cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode))
	if err != nil {
		logger.Fatal().
			Str("app", "failed to connect db").
			Msg(err.Error())
		return
	}

	logger.Debug().
		Str("cfg.DBHost", cfg.DBHost).
		Str("cfg.DBUser", cfg.DBUser).
		Str("cfg.DBName", cfg.DBName).
		Str("cfg.DBPassword", cfg.DBPassword).
		Str("cfg.DBSSLMode", cfg.DBSSLMode).
		Msg("db connected")

	repo := pg.NewPgRepository(db, logger)
	cache := cache.NewChatCache(service.Options().Store)
	serv := NewChatService(repo, logger, flowClient, botClient, cache, service.Options().Broker)

	if err := pb.RegisterChatServiceHandler(service.Server(), serv); err != nil {
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
