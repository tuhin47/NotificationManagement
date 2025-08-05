package cmd

import (
	"NotificationManagement/config"
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

	// TODO: Initialize connections when conn package is available
	// conn.ConnectDb()
	// conn.ConnectRedis()
	// conn.InitAsynqClient()
	// conn.InitAsyncInspector()
	// conn.ConnectEmail()

	// asynq connections
	// conn.InitAsynqClient()
	// conn.InitAsyncInspector()

	if err := RootCmd.Execute(); err != nil {
		logger.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
