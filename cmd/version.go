package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "get version of chat server",
	Long:  "get version of chat server",
	Run:   getVersion,
}

func getVersion(cmd *cobra.Command, _ []string) {
	fmt.Printf("version: %s\n", version)
}
