package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// The startCmd lets us start our application.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start chat server",
	Long:  "start chat server",
	Run:   app,
}

// app starts the application
func app(cmd *cobra.Command, _ []string) {
	cfg := new(Config)
	err := NewConfig(cmd.Flags(), "config", cfg)
	if err != nil {
		log.Fatalf("app: failed to load config %#v", err)
	}
	app, err := NewApp(cfg)
	if err != nil {
		log.Fatalf("app: failed to create app %#v", err)
	}
	run(app.Start, app.Stop)
}

func run(start, stop func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	complete := make(chan bool, 1)
	go func() {
		start()
		complete <- true
	}()
	select {
	case <-signals:
		stop()
	case <-complete:
	}
}
