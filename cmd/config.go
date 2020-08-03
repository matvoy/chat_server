package main

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	configTag = "config"
)

type Configurator interface {
	Env() string
}

type Config struct {
	Environment      string `config:"env"`
	AppPort          int    `config:"app_port"`
	Debug            bool   `config:"debug"`
	TelegramBotToken string `config:"telegram_bot_token"`
}

func (c *Config) Env() string {
	return c.Environment
}

func NewConfig(flags *pflag.FlagSet, fileName string, c Configurator) error {
	vi := viper.New()
	vi.AutomaticEnv()
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// extract the filename to use in building paths
	_, file := path.Split(os.Args[0])
	// set up File loading
	vi.SetConfigName(fileName)
	vi.AddConfigPath(fmt.Sprintf("/etc/%s/", file))
	vi.AddConfigPath(fmt.Sprintf("$HOME/.%s", file))
	vi.AddConfigPath(".")
	// bind the command line flags
	if flags != nil {
		if err := vi.BindPFlags(flags); err != nil {
			// TO DO LOG ERROR
		}
	}
	// read from the config file
	err := vi.ReadInConfig()
	if err == nil && vi.ConfigFileUsed() != "" {
		// TO DO LOG CONFIGURATION LOADED
	}
	t := reflect.TypeOf(c)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// bind the possible environment variables
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if err = vi.BindEnv(field.Tag.Get(configTag)); err != nil {
			// TO DO LOG ERROR
		}
	}
	// unmarshal configs into the provided struct
	if err = vi.Unmarshal(&c, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = configTag
	}); err != nil {
		return err
	}
	return nil
}
