package cmd

import (
	aws "NotificationManagement/aws-config"
	"NotificationManagement/config"
	"context"
	"fmt"

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
			fmt.Printf("Error creating AWS Config client: %v\n", err)
			return
		}

		ctx := context.Background()
		rules, err := client.ListConfigRules(ctx)
		if err != nil {
			fmt.Printf("Error listing config rules: %v\n", err)
			return
		}

		fmt.Printf("Found %d config rules:\n", len(rules))
		for _, rule := range rules {
			fmt.Printf("- %s: %s\n", *rule.ConfigRuleName, *rule.Description)
		}
	},
}

var checkStatusCmd = &cobra.Command{
	Use:   "check-status",
	Short: "Check AWS Config recorder status",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := aws.NewConfigClient()
		if err != nil {
			fmt.Printf("Error creating AWS Config client: %v\n", err)
			return
		}

		ctx := context.Background()
		status, err := client.GetConfigurationRecorderStatus(ctx)
		if err != nil {
			fmt.Printf("Error getting recorder status: %v\n", err)
			return
		}

		fmt.Printf("Configuration Recorder Status:\n")
		fmt.Printf("- Name: %s\n", *status.Name)
		fmt.Printf("- Recording: %v\n", status.Recording)
		fmt.Printf("- Last Status: %s\n", status.LastStatus)
		if status.LastErrorCode != nil {
			fmt.Printf("- Last Error: %s - %s\n", *status.LastErrorCode, *status.LastErrorMessage)
		}
	},
}

var testConnectionCmd = &cobra.Command{
	Use:   "test-connection",
	Short: "Test AWS Config service connection",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Testing AWS Config service connection...\n")
		fmt.Printf("Region: %s\n", config.AWS().Region)
		fmt.Printf("Using LocalStack: %v\n", config.AWS().UseLocalStack)
		if config.AWS().UseLocalStack {
			fmt.Printf("LocalStack Endpoint: %s\n", config.AWS().Endpoint)
		}

		client, err := aws.NewConfigClient()
		if err != nil {
			fmt.Printf("Error creating AWS Config client: %v\n", err)
			return
		}

		fmt.Println("✅ AWS Config client created successfully!")

		// Try to list rules to test the connection
		ctx := context.Background()
		_, err = client.ListConfigRules(ctx)
		if err != nil {
			fmt.Printf("❌ Error listing config rules: %v\n", err)
			return
		}

		fmt.Println("✅ AWS Config service connection successful!")
	},
}

func init() {
	configCmd.AddCommand(listRulesCmd)
	configCmd.AddCommand(checkStatusCmd)
	configCmd.AddCommand(testConnectionCmd)
	RootCmd.AddCommand(configCmd)
}
