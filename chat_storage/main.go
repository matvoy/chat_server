package main

import (
	"database/sql"
	"fmt"
	"os"

	pb "github.com/matvoy/chat_server/chat_storage/proto/storage"
	"github.com/matvoy/chat_server/chat_storage/repo/pg"

	_ "github.com/lib/pq"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
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

func main() {
	cfg := &Config{}
	service := micro.NewService(
		micro.Name("webitel.chat.service.storage"),
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
			return nil
		}),
	)
	logger, err := NewLogger(cfg.LogLevel)
	if err != nil {
		logger.Fatal().
			Str("app", "failed to parse log level").
			Msg(err.Error())
		return
	}
	logger.Debug().
		Msg("logger created")
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
	serv := NewStorageService(repo, logger)

	if err := pb.RegisterStorageServiceHandler(service.Server(), serv); err != nil {
		logger.Fatal().
			Str("app", "failed to register service").
			Msg(err.Error())
		return
	}
	logger.Debug().
		Msg("service handler registered")
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
