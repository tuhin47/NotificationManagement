package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	configtype "github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"os"
)

type AWSClient struct {
	config    *configservice.Client
	ssm       *ssm.Client
	awsConfig *AWSConfig
}

func NewAWSClient(cnf *AWSConfig) (*AWSClient, error) {
	var creds aws.CredentialsProvider
	if cnf.AccessKeyID != "" && cnf.SecretAccessKey != "" {
		creds = credentials.NewStaticCredentialsProvider(
			cnf.AccessKeyID,
			cnf.SecretAccessKey,
			"",
		)
	}
	var opts []func(*awsconfig.LoadOptions) error

	if cnf.Region != "" {
		opts = append(opts, awsconfig.WithRegion(cnf.Region))
	}

	if creds != nil {
		opts = append(opts, awsconfig.WithCredentialsProvider(creds))
	}

	// Load AWS awsconfig
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create SSM client with custom endpoint
	configOptions := []func(*configservice.Options){
		func(o *configservice.Options) {
			if cnf.Endpoint != "" {
				o.BaseEndpoint = aws.String(cnf.Endpoint)
			}
		},
	}
	ssmOptions := []func(*ssm.Options){
		func(o *ssm.Options) {
			if cnf.Endpoint != "" {
				o.BaseEndpoint = aws.String(cnf.Endpoint)
			}
		},
	}

	return &AWSClient{
		config:    configservice.NewFromConfig(cfg, configOptions...),
		ssm:       ssm.NewFromConfig(cfg, ssmOptions...),
		awsConfig: cnf,
	}, nil
}

func (c *AWSClient) ListConfigRules(ctx context.Context) ([]configtype.ConfigRule, error) {
	input := &configservice.DescribeConfigRulesInput{}

	result, err := c.config.DescribeConfigRules(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe config rules: %v", err)
	}

	return result.ConfigRules, nil
}

func (c *AWSClient) GetConfigurationRecorderStatus(ctx context.Context) (*configtype.ConfigurationRecorderStatus, error) {
	input := &configservice.DescribeConfigurationRecorderStatusInput{}

	result, err := c.config.DescribeConfigurationRecorderStatus(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe configuration recorder status: %v", err)
	}

	if len(result.ConfigurationRecordersStatus) == 0 {
		return nil, fmt.Errorf("no configuration recorder found")
	}

	return &result.ConfigurationRecordersStatus[0], nil
}

func (c *AWSClient) loadFromSsm() {
	if os.Getenv(EnvConfigFromSSM) == "false" {
		return
	}
	resp, err := c.ssm.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           &c.awsConfig.ConfigService.SSM,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		panic("failed to get config from SSM: " + err.Error())
	}
	var ssmConfigMap Config
	err = json.Unmarshal([]byte(*resp.Parameter.Value), &ssmConfigMap)
	if err != nil {
		panic("failed to unmarshal config from SSM: " + err.Error())
	}

	setViperFields(ssmConfigMap, "")
}
