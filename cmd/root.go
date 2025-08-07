package cmd

import (
	"NotificationManagement/config"
	"NotificationManagement/conn"
	"NotificationManagement/logger"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use: "app",
}

func init() {
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(workerCmd)
}

// Execute executes the root command
func Execute() {
	// load config
	config.LoadConfig()

	// Initialize logger
	logger.Init()
	defer logger.Sync()

	conn.ConnectRedis()
	// conn.ConnectEmail()

	if err := RootCmd.Execute(); err != nil {
		logger.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
