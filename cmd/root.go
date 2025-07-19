package cmd

import (
	"fmt"
	"os"

	"NotificationManagement/config"

	"github.com/spf13/cobra"
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
	initLogger()

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
		fmt.Println(err)
		os.Exit(1)
	}
}

func initLogger() {
	fmt.Println("Initializing logger...")
	fmt.Println("Logger file path:", config.Logger().FilePath)
	// TODO: Initialize logger when logger package is available
	// logger.SetFileLogger(config.Logger().FilePath)
}
