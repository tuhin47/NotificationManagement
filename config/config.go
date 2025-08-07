package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/viper"
)

type Config struct {
	App         AppConfig      `mapstructure:"app"`
	Database    DatabaseConfig `mapstructure:"database"`
	AsynqConfig AsynqConfig    `mapstructure:"asynq"`
	Redis       RedisConfig    `mapstructure:"redis"`
	Email       EmailConfig    `mapstructure:"email"`
	AWS         AWSConfig      `mapstructure:"aws"`
	Logger      LoggerConfig   `mapstructure:"logger"`
	Keycloak    KeycloakConfig `mapstructure:"keycloak"`
}

type AppConfig struct {
	Name       string `mapstructure:"name"`
	Version    string `mapstructure:"version"`
	Port       int    `mapstructure:"port"`
	Env        string `mapstructure:"env"`
	Encryption string `mapstructure:"encryption"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type AsynqConfig struct {
	RedisAddr                        string        `mapstructure:"redisAddr"`
	DB                               int           `mapstructure:"db"`
	Pass                             string        `mapstructure:"pass"`
	Concurrency                      int           `mapstructure:"concurrency"`
	Queue                            string        `mapstructure:"queue"`
	Retention                        time.Duration `mapstructure:"retention"` // in hours
	RetryCount                       int           `mapstructure:"retryCount"`
	Delay                            time.Duration `mapstructure:"delay"` // in seconds
	EmailInvitationTaskDelay         time.Duration `mapstructure:"emailInvitationTaskDelay"`
	EmailInvitationTaskRetryCount    int           `mapstructure:"emailInvitationTaskRetryCount"`
	EmailInvitationTaskRetryDelay    time.Duration `mapstructure:"emailInvitationTaskRetryDelay"`
	EventReminderTaskRetryCount      int           `mapstructure:"eventReminderTaskRetryCount"`
	EventReminderTaskRetryDelay      time.Duration `mapstructure:"eventReminderTaskRetryDelay"`
	EventReminderEmailTaskDelay      time.Duration `mapstructure:"eventReminderEmailTaskDelay"`
	EventReminderEmailTaskRetryCount int           `mapstructure:"eventReminderEmailTaskRetryCount"`
	EventReminderEmailTaskRetryDelay time.Duration `mapstructure:"eventReminderEmailTaskRetryDelay"`
}

type RedisConfig struct {
	Host               string        `mapstructure:"host"`
	Port               int           `mapstructure:"port"`
	Password           string        `mapstructure:"password"`
	DB                 int           `mapstructure:"db"`
	MandatoryPrefix    string        `mapstructure:"mandatoryPrefix"`
	AccessUuidPrefix   string        `mapstructure:"accessUuidPrefix"`
	RefreshUuidPrefix  string        `mapstructure:"refreshUuidPrefix"`
	UserPrefix         string        `mapstructure:"userPrefix"`
	PermissionPrefix   string        `mapstructure:"permissionPrefix"`
	UserCacheTTL       time.Duration `mapstructure:"userCacheTTL"`
	PermissionCacheTTL time.Duration `mapstructure:"permissionCacheTTL"`
}

type EmailConfig struct {
	Url      string
	Timeout  time.Duration
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

type KeycloakConfig struct {
	ServerURL    string `mapstructure:"server_url"`
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type ConfigServiceConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	SSM     string `mapstructure:"ssm"`

	// Add specific config service settings as needed
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
	Format   string `mapstructure:"format"`
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
	EnvLogFormat   = "LOGGER_FORMAT"

	// Keycloak environment variables
	EnvKeycloakServerURL    = "KEYCLOAK_SERVER_URL"
	EnvKeycloakRealm        = "KEYCLOAK_REALM"
	EnvKeycloakClientID     = "KEYCLOAK_CLIENT_ID"
	EnvKeycloakClientSecret = "KEYCLOAK_CLIENT_SECRET"

	EnvAPIKeyEncryptionSecret = "API_KEY_ENCRYPTION_SECRET"
)

// LoadConfig loads configuration from file, environment variables, or SSM Parameter Store
func LoadConfig() {
	// Set default values
	defaultConfig := getDefaults()
	viper.SetDefault("app", defaultConfig.App)
	viper.SetDefault("database", defaultConfig.Database)
	viper.SetDefault("redis", defaultConfig.Redis)
	viper.SetDefault("email", defaultConfig.Email)
	viper.SetDefault("aws", defaultConfig.AWS)
	viper.SetDefault("logger", defaultConfig.Logger)
	viper.SetDefault("keycloak", defaultConfig.Keycloak)

	if os.Getenv(EnvConfigFromSSM) != "false" {
		ssmParam := os.Getenv(EnvConfigSSMParam)
		if ssmParam == "" {
			ssmParam = defaultConfig.AWS.ConfigService.SSM
		}
		region := os.Getenv(EnvAWSRegion)
		if region == "" {
			region = defaultConfig.AWS.Region
		}
		endpoint := os.Getenv(EnvAWSEndpoint)
		if endpoint == "" {
			endpoint = defaultConfig.AWS.Endpoint
		}
		useLocalEnv := os.Getenv(EnvAWSUseLocalStack)
		useLocalStack := useLocalEnv == "true"
		if useLocalEnv == "" {
			useLocalStack = defaultConfig.AWS.UseLocalStack
		}

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
		var ssmConfigMap map[string]interface{}
		err = json.Unmarshal([]byte(*resp.Parameter.Value), &ssmConfigMap)
		if err != nil {
			panic("failed to unmarshal config from SSM: " + err.Error())
		}
		if err := viper.MergeConfigMap(ssmConfigMap); err != nil {
			panic("failed to merge SSM config: " + err.Error())
		}
	}

	loadFromEnv()

	// Unmarshal config
	if err := viper.Unmarshal(&appConfig); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config Loaded")
}

func getDefaults() *Config {
	conf := &Config{}
	conf.App = AppConfig{
		Name:       "NotificationManagement",
		Version:    "1.0.0",
		Port:       8080,
		Env:        "development",
		Encryption: "laeoGcA0ZFFsm3d9SUKevwG4VL4QN9Yi",
	}
	conf.Database = DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		Name:     "notification_management",
		SSLMode:  "disable",
	}
	conf.Redis = RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}
	conf.Email = EmailConfig{
		Host:     "localhost",
		Port:     587,
		Username: "",
		Password: "",
		From:     "noreply@example.com",
	}
	conf.AWS = AWSConfig{
		Region:          "us-east-1",
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		Endpoint:        "http://localhost:4566",
		UseLocalStack:   true,
		ConfigService: ConfigServiceConfig{
			Enabled: true,
			SSM:     "/myapp/config",
		},
	}
	conf.Logger = LoggerConfig{
		Level:    "info",
		FilePath: "logs/app.log",
		Format:   "console",
	}
	conf.Keycloak = KeycloakConfig{
		ServerURL:    "http://localhost:8081",
		Realm:        "gocloak",
		ClientID:     "gocloak",
		ClientSecret: "gocloak-secret",
	}
	return conf
}

func loadFromEnv() {
	// App environment variables
	_ = viper.BindEnv("app.name", EnvAppName)
	_ = viper.BindEnv("app.version", EnvAppVersion)
	_ = viper.BindEnv("app.port", EnvAppPort)
	_ = viper.BindEnv("app.env", EnvAppEnv)
	_ = viper.BindEnv("app.encryption", EnvAPIKeyEncryptionSecret)

	// Database environment variables
	_ = viper.BindEnv("database.host", EnvDBHost)
	_ = viper.BindEnv("database.port", EnvDBPort)
	_ = viper.BindEnv("database.user", EnvDBUser)
	_ = viper.BindEnv("database.password", EnvDBPassword)
	_ = viper.BindEnv("database.name", EnvDBName)
	_ = viper.BindEnv("database.ssl_mode", EnvDBSSLMode)

	// Redis environment variables
	_ = viper.BindEnv("redis.host", EnvRedisHost)
	_ = viper.BindEnv("redis.port", EnvRedisPort)
	_ = viper.BindEnv("redis.password", EnvRedisPassword)
	_ = viper.BindEnv("redis.db", EnvRedisDB)

	// Email environment variables
	_ = viper.BindEnv("email.host", EnvEmailHost)
	_ = viper.BindEnv("email.port", EnvEmailPort)
	_ = viper.BindEnv("email.username", EnvEmailUsername)
	_ = viper.BindEnv("email.password", EnvEmailPassword)
	_ = viper.BindEnv("email.from", EnvEmailFrom)

	// AWS environment variables
	_ = viper.BindEnv("aws.region", EnvAWSRegion)
	_ = viper.BindEnv("aws.region", EnvAWSRegion)
	_ = viper.BindEnv("aws.region", EnvAWSRegion)
	_ = viper.BindEnv("aws.access_key_id", EnvAWSAccessKeyID)
	_ = viper.BindEnv("aws.secret_access_key", EnvAWSSecretAccessKey)
	_ = viper.BindEnv("aws.endpoint", EnvAWSEndpoint)
	_ = viper.BindEnv("aws.use_localstack", EnvAWSUseLocalStack)
	_ = viper.BindEnv("aws.config_service.enabled", EnvAWSConfigServiceEnabled)

	// Logger environment variables
	_ = viper.BindEnv("logger.level", EnvLogLevel)
	_ = viper.BindEnv("logger.file_path", EnvLogFilePath)
	_ = viper.BindEnv("logger.format", EnvLogFormat)

	// Keycloak environment variables
	_ = viper.BindEnv("keycloak.server_url", EnvKeycloakServerURL)
	_ = viper.BindEnv("keycloak.realm", EnvKeycloakRealm)
	_ = viper.BindEnv("keycloak.client_id", EnvKeycloakClientID)
	_ = viper.BindEnv("keycloak.client_secret", EnvKeycloakClientSecret)
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

func Asynq() AsynqConfig {
	return appConfig.AsynqConfig
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

// Keycloak returns the keycloak configuration
func Keycloak() KeycloakConfig {
	return appConfig.Keycloak
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
