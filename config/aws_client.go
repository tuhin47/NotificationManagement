package config

import (
	"context"
	"fmt"

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
	appConfig := AWS()

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
