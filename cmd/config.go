package cmd

import (
	aws "NotificationManagement/aws-config"
	"NotificationManagement/config"
	"NotificationManagement/logger"
	"context"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "AWS Config service operations",
	Long:  `Perform AWS Config service operations like listing rules, checking compliance, etc.`,
}

var listRulesCmd = &cobra.Command{
	Use:   "list-rules",
	Short: "List all AWS Config rules",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := aws.NewConfigClient()
		if err != nil {
			logger.Error("Error creating AWS Config client", "error", err)
			return
		}

		ctx := context.Background()
		rules, err := client.ListConfigRules(ctx)
		if err != nil {
			logger.Error("Error listing config rules", "error", err)
			return
		}

		logger.Info("Found config rules", "count", len(rules))
		for _, rule := range rules {
			logger.Info("Config rule", "name", *rule.ConfigRuleName, "description", *rule.Description)
		}
	},
}

var checkStatusCmd = &cobra.Command{
	Use:   "check-status",
	Short: "Check AWS Config recorder status",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := aws.NewConfigClient()
		if err != nil {
			logger.Error("Error creating AWS Config client", "error", err)
			return
		}

		ctx := context.Background()
		status, err := client.GetConfigurationRecorderStatus(ctx)
		if err != nil {
			logger.Error("Error getting recorder status", "error", err)
			return
		}

		logger.Info("Configuration Recorder Status")
		logger.Info("Name", "name", *status.Name)
		logger.Info("Recording", "recording", status.Recording)
		logger.Info("Last Status", "last_status", string(status.LastStatus))
		if status.LastErrorCode != nil {
			logger.Info("Last Error", "error_code", *status.LastErrorCode, "error_message", *status.LastErrorMessage)
		}
	},
}

var testConnectionCmd = &cobra.Command{
	Use:   "test-connection",
	Short: "Test AWS Config service connection",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Testing AWS Config service connection...")
		logger.Info("Region", "region", config.AWS().Region)
		logger.Info("Using LocalStack", "use_local_stack", config.AWS().UseLocalStack)
		if config.AWS().UseLocalStack {
			logger.Info("LocalStack Endpoint", "endpoint", config.AWS().Endpoint)
		}

		client, err := aws.NewConfigClient()
		if err != nil {
			logger.Error("Error creating AWS Config client", "error", err)
			return
		}

		logger.Info("✅ AWS Config client created successfully!")

		// Try to list rules to test the connection
		ctx := context.Background()
		_, err = client.ListConfigRules(ctx)
		if err != nil {
			logger.Error("❌ Error listing config rules", "error", err)
			return
		}

		logger.Info("✅ AWS Config service connection successful!")
	},
}

func init() {
	configCmd.AddCommand(listRulesCmd)
	configCmd.AddCommand(checkStatusCmd)
	configCmd.AddCommand(testConnectionCmd)
	RootCmd.AddCommand(configCmd)
}
