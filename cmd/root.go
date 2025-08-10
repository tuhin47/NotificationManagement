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

func Execute() {
	config.LoadConfig()

	logger.Init()
	defer logger.Sync()

	conn.ConnectRedis()

	if err := RootCmd.Execute(); err != nil {
		logger.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
