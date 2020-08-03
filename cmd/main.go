package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chat_server",
	Short: "Webitel Chat Server",
	Long:  "Webitel Chat Server",
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flags := startCmd.Flags()
	flags.Int("app_port", 3000, "http port to listen to traffic on")
	flags.String("debug", "false", "debug log level")
	flags.String("telegram_bot_token", "token", "telegram bot token")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err.Error())
		os.Exit(2)
	}
}
