package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Email    EmailConfig    `mapstructure:"email"`
	AWS      AWSConfig      `mapstructure:"aws"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
	Env     string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type AWSConfig struct {
	Region          string              `mapstructure:"region"`
	AccessKeyID     string              `mapstructure:"access_key_id"`
	SecretAccessKey string              `mapstructure:"secret_access_key"`
	Endpoint        string              `mapstructure:"endpoint"`
	UseLocalStack   bool                `mapstructure:"use_localstack"`
	ConfigService   ConfigServiceConfig `mapstructure:"config_service"`
}

type ConfigServiceConfig struct {
	Enabled bool `mapstructure:"enabled"`
	// Add specific config service settings as needed
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
}

var (
	appConfig *Config
)

const (
	EnvConfigFromSSM      = "CONFIG_FROM_SSM"
	EnvConfigSSMParam     = "CONFIG_SSM_PARAM"
	EnvAWSRegion          = "AWS_REGION"
	EnvAWSEndpoint        = "AWS_ENDPOINT"
	EnvAWSUseLocalStack   = "AWS_USE_LOCALSTACK"
	EnvAWSAccessKeyID     = "AWS_ACCESS_KEY_ID"
	EnvAWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	EnvAppName    = "APP_NAME"
	EnvAppVersion = "APP_VERSION"
	EnvAppPort    = "APP_PORT"
	EnvAppEnv     = "APP_ENV"

	EnvDBHost     = "DB_HOST"
	EnvDBPort     = "DB_PORT"
	EnvDBUser     = "DB_USER"
	EnvDBPassword = "DB_PASSWORD"
	EnvDBName     = "DB_NAME"
	EnvDBSSLMode  = "DB_SSL_MODE"

	EnvRedisHost     = "REDIS_HOST"
	EnvRedisPort     = "REDIS_PORT"
	EnvRedisPassword = "REDIS_PASSWORD"
	EnvRedisDB       = "REDIS_DB"

	EnvEmailHost     = "EMAIL_HOST"
	EnvEmailPort     = "EMAIL_PORT"
	EnvEmailUsername = "EMAIL_USERNAME"
	EnvEmailPassword = "EMAIL_PASSWORD"
	EnvEmailFrom     = "EMAIL_FROM"

	EnvAWSConfigServiceEnabled = "AWS_CONFIG_SERVICE_ENABLED"

	EnvLogLevel    = "LOG_LEVEL"
	EnvLogFilePath = "LOG_FILE_PATH"
)

// LoadConfig loads configuration from file, environment variables, or SSM Parameter Store
func LoadConfig() {
	if os.Getenv(EnvConfigFromSSM) == "true" {
		ssmParam := os.Getenv(EnvConfigSSMParam)
		if ssmParam == "" {
			ssmParam = "/myapp/config"
		}
		region := os.Getenv(EnvAWSRegion)
		if region == "" {
			region = "us-east-1"
		}
		endpoint := os.Getenv(EnvAWSEndpoint)
		useLocalStack := os.Getenv(EnvAWSUseLocalStack) == "true"

		var awsCfg aws.Config
		var err error
		if useLocalStack {
			customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			})
			awsCfg, err = config.LoadDefaultConfig(context.TODO(),
				config.WithEndpointResolverWithOptions(customResolver),
				config.WithRegion(region),
				config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
					return aws.Credentials{
						AccessKeyID:     os.Getenv(EnvAWSAccessKeyID),
						SecretAccessKey: os.Getenv(EnvAWSSecretAccessKey),
					}, nil
				})),
			)
		} else {
			awsCfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		}
		if err != nil {
			panic("failed to load AWS config: " + err.Error())
		}
		ssmClient := ssm.NewFromConfig(awsCfg)
		resp, err := ssmClient.GetParameter(context.TODO(), &ssm.GetParameterInput{
			Name:           &ssmParam,
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			panic("failed to get config from SSM: " + err.Error())
		}
		var loadedConfig Config
		err = json.Unmarshal([]byte(*resp.Parameter.Value), &loadedConfig)
		if err != nil {
			panic("failed to unmarshal config from SSM: " + err.Error())
		}
		appConfig = &loadedConfig
		return
	}
	// Set default values
	setDefaults()

	// Read from environment variables
	loadFromEnv()

	// Unmarshal config
	if err := viper.Unmarshal(&appConfig); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		os.Exit(1)
	}
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "NotificationManagement")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.env", "development")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "notification_management")
	viper.SetDefault("database.ssl_mode", "disable")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Email defaults
	viper.SetDefault("email.host", "localhost")
	viper.SetDefault("email.port", 587)
	viper.SetDefault("email.username", "")
	viper.SetDefault("email.password", "")
	viper.SetDefault("email.from", "noreply@example.com")

	// AWS defaults
	viper.SetDefault("aws.region", "us-east-1")
	viper.SetDefault("aws.access_key_id", "test")
	viper.SetDefault("aws.secret_access_key", "test")
	viper.SetDefault("aws.endpoint", "")
	viper.SetDefault("aws.use_localstack", false)
	viper.SetDefault("aws.config_service.enabled", true)

	// Logger defaults
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.file_path", "logs/app.log")
}

func loadFromEnv() {
	// App environment variables
	viper.BindEnv("app.name", EnvAppName)
	viper.BindEnv("app.version", EnvAppVersion)
	viper.BindEnv("app.port", EnvAppPort)
	viper.BindEnv("app.env", EnvAppEnv)

	// Database environment variables
	viper.BindEnv("database.host", EnvDBHost)
	viper.BindEnv("database.port", EnvDBPort)
	viper.BindEnv("database.user", EnvDBUser)
	viper.BindEnv("database.password", EnvDBPassword)
	viper.BindEnv("database.name", EnvDBName)
	viper.BindEnv("database.ssl_mode", EnvDBSSLMode)

	// Redis environment variables
	viper.BindEnv("redis.host", EnvRedisHost)
	viper.BindEnv("redis.port", EnvRedisPort)
	viper.BindEnv("redis.password", EnvRedisPassword)
	viper.BindEnv("redis.db", EnvRedisDB)

	// Email environment variables
	viper.BindEnv("email.host", EnvEmailHost)
	viper.BindEnv("email.port", EnvEmailPort)
	viper.BindEnv("email.username", EnvEmailUsername)
	viper.BindEnv("email.password", EnvEmailPassword)
	viper.BindEnv("email.from", EnvEmailFrom)

	// AWS environment variables
	viper.BindEnv("aws.region", EnvAWSRegion)
	viper.BindEnv("aws.access_key_id", EnvAWSAccessKeyID)
	viper.BindEnv("aws.secret_access_key", EnvAWSSecretAccessKey)
	viper.BindEnv("aws.endpoint", EnvAWSEndpoint)
	viper.BindEnv("aws.use_localstack", EnvAWSUseLocalStack)
	viper.BindEnv("aws.config_service.enabled", EnvAWSConfigServiceEnabled)

	// Logger environment variables
	viper.BindEnv("logger.level", EnvLogLevel)
	viper.BindEnv("logger.file_path", EnvLogFilePath)
}

// GetConfig returns the application configuration
func GetConfig() *Config {
	return appConfig
}

// App returns the app configuration
func App() AppConfig {
	return appConfig.App
}

// Database returns the database configuration
func Database() DatabaseConfig {
	return appConfig.Database
}

// Redis returns the redis configuration
func Redis() RedisConfig {
	return appConfig.Redis
}

// Email returns the email configuration
func Email() EmailConfig {
	return appConfig.Email
}

// AWS returns the AWS configuration
func AWS() AWSConfig {
	return appConfig.AWS
}

// Logger returns the logger configuration
func Logger() LoggerConfig {
	return appConfig.Logger
}

// GetDSN returns the database connection string
func GetDSN() string {
	db := Database()
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode)
}

// GetRedisAddr returns the Redis connection address
func GetRedisAddr() string {
	redis := Redis()
	return fmt.Sprintf("%s:%d", redis.Host, redis.Port)
}

// IsDevelopment returns true if the application is running in development mode
func IsDevelopment() bool {
	return App().Env == "development"
}

// IsProduction returns true if the application is running in production mode
func IsProduction() bool {
	return App().Env == "production"
}

// IsLocalStack returns true if LocalStack should be used
func IsLocalStack() bool {
	return AWS().UseLocalStack
}
