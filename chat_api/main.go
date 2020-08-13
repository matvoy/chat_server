package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chat_api",
	Short: "Webitel Chat API",
	Long:  "Webitel Chat API",
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flags := startCmd.Flags()
	flags.Int("api_http_port", 3000, "http port")
	flags.String("log_level", "debug", "debug log level")
	flags.String("db_host", "localhost", "the database host")
	flags.String("db_user", "postgres", "the database user")
	flags.String("db_name", "postgres", "the database name")
	flags.String("db_password", "postgres", "the database password")
	flags.String("db_sslmode", "disable", "the database ssl mode")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err.Error())
		os.Exit(2)
	}
}
