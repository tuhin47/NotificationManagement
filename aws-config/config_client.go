package aws_config

import (
	"context"
	"fmt"
	"log"

	"NotificationManagement/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
)

type ConfigClient struct {
	client *configservice.Client
	cfg    aws.Config
}

// NewConfigClient creates a new AWS Config service client
func NewConfigClient() (*ConfigClient, error) {
	appConfig := config.AWS()

	var cfg aws.Config
	var err error

	if appConfig.UseLocalStack {
		// Use LocalStack configuration
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           appConfig.Endpoint,
				SigningRegion: appConfig.Region,
			}, nil
		})

		cfg, err = awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithEndpointResolverWithOptions(customResolver),
			awsconfig.WithRegion(appConfig.Region),
			awsconfig.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     appConfig.AccessKeyID,
					SecretAccessKey: appConfig.SecretAccessKey,
				}, nil
			})),
		)
	} else {
		// Use real AWS configuration
		cfg, err = awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithRegion(appConfig.Region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := configservice.NewFromConfig(cfg)

	return &ConfigClient{
		client: client,
		cfg:    cfg,
	}, nil
}

// ListConfigRules lists all Config rules
func (c *ConfigClient) ListConfigRules(ctx context.Context) ([]types.ConfigRule, error) {
	input := &configservice.DescribeConfigRulesInput{}

	result, err := c.client.DescribeConfigRules(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe config rules: %v", err)
	}

	return result.ConfigRules, nil
}

// GetConfigRule gets a specific Config rule by name
func (c *ConfigClient) GetConfigRule(ctx context.Context, ruleName string) (*types.ConfigRule, error) {
	input := &configservice.DescribeConfigRulesInput{
		ConfigRuleNames: []string{ruleName},
	}

	result, err := c.client.DescribeConfigRules(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe config rule %s: %v", ruleName, err)
	}

	if len(result.ConfigRules) == 0 {
		return nil, fmt.Errorf("config rule %s not found", ruleName)
	}

	return &result.ConfigRules[0], nil
}

// CreateConfigRule creates a new Config rule
func (c *ConfigClient) CreateConfigRule(ctx context.Context, ruleName, description string) error {
	input := &configservice.PutConfigRuleInput{
		ConfigRule: &types.ConfigRule{
			ConfigRuleName: aws.String(ruleName),
			Description:    aws.String(description),
			// TODO: Add proper scope and source configuration
		},
	}

	_, err := c.client.PutConfigRule(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create config rule %s: %v", ruleName, err)
	}

	log.Printf("Successfully created config rule: %s", ruleName)
	return nil
}

// DeleteConfigRule deletes a Config rule
func (c *ConfigClient) DeleteConfigRule(ctx context.Context, ruleName string) error {
	input := &configservice.DeleteConfigRuleInput{
		ConfigRuleName: aws.String(ruleName),
	}

	_, err := c.client.DeleteConfigRule(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete config rule %s: %v", ruleName, err)
	}

	log.Printf("Successfully deleted config rule: %s", ruleName)
	return nil
}

// GetComplianceDetails gets compliance details for resources
func (c *ConfigClient) GetComplianceDetails(ctx context.Context, resourceType, resourceId string) ([]types.EvaluationResult, error) {
	input := &configservice.GetComplianceDetailsByConfigRuleInput{
		ConfigRuleName: aws.String(resourceType),
		ComplianceTypes: []types.ComplianceType{
			types.ComplianceTypeCompliant,
			types.ComplianceTypeNonCompliant,
		},
	}

	result, err := c.client.GetComplianceDetailsByConfigRule(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get compliance details: %v", err)
	}

	return result.EvaluationResults, nil
}

// StartConfigRulesEvaluation starts evaluation for Config rules
func (c *ConfigClient) StartConfigRulesEvaluation(ctx context.Context, ruleNames []string) error {
	input := &configservice.StartConfigRulesEvaluationInput{
		ConfigRuleNames: ruleNames,
	}

	_, err := c.client.StartConfigRulesEvaluation(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to start config rules evaluation: %v", err)
	}

	log.Printf("Started evaluation for config rules: %v", ruleNames)
	return nil
}

// GetConfigurationRecorderStatus gets the status of configuration recorder
func (c *ConfigClient) GetConfigurationRecorderStatus(ctx context.Context) (*types.ConfigurationRecorderStatus, error) {
	input := &configservice.DescribeConfigurationRecorderStatusInput{}

	result, err := c.client.DescribeConfigurationRecorderStatus(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe configuration recorder status: %v", err)
	}

	if len(result.ConfigurationRecordersStatus) == 0 {
		return nil, fmt.Errorf("no configuration recorder found")
	}

	return &result.ConfigurationRecordersStatus[0], nil
}

// StartConfigurationRecorder starts the configuration recorder
func (c *ConfigClient) StartConfigurationRecorder(ctx context.Context, recorderName string) error {
	input := &configservice.StartConfigurationRecorderInput{
		ConfigurationRecorderName: aws.String(recorderName),
	}

	_, err := c.client.StartConfigurationRecorder(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to start configuration recorder %s: %v", recorderName, err)
	}

	log.Printf("Successfully started configuration recorder: %s", recorderName)
	return nil
}

// StopConfigurationRecorder stops the configuration recorder
func (c *ConfigClient) StopConfigurationRecorder(ctx context.Context, recorderName string) error {
	input := &configservice.StopConfigurationRecorderInput{
		ConfigurationRecorderName: aws.String(recorderName),
	}

	_, err := c.client.StopConfigurationRecorder(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to stop configuration recorder %s: %v", recorderName, err)
	}

	log.Printf("Successfully stopped configuration recorder: %s", recorderName)
	return nil
}
