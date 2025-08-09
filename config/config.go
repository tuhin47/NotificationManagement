package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/viper"
)

var properties = []string{"app", "database", "redis", "asynq", "email", "telegram", "aws", "keycloak", "logger", "config_service", "development"}

type Config struct {
	App         AppConfig         `mapstructure:"app"`
	Database    DatabaseConfig    `mapstructure:"database"`
	AsynqConfig AsynqConfig       `mapstructure:"asynq" json:"asynq"`
	Redis       RedisConfig       `mapstructure:"redis"`
	Email       EmailConfig       `mapstructure:"email"`
	AWS         AWSConfig         `mapstructure:"aws"`
	Logger      LoggerConfig      `mapstructure:"logger"`
	Keycloak    KeycloakConfig    `mapstructure:"keycloak"`
	Telegram    TelegramConfig    `mapstructure:"telegram"`
	Development DevelopmentConfig `mapstructure:"development"`
}
type DevelopmentConfig struct {
	GeminiKey string `mapstructure:"geminikey"`
}

type AppConfig struct {
	Name       string `mapstructure:"name"`
	Version    string `mapstructure:"version"`
	Port       *int   `mapstructure:"port"`
	Env        string `mapstructure:"env"`
	Encryption string `mapstructure:"encryption"`
	Domain     string `mapstructure:"domain"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     *int   `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type AsynqConfig struct {
	RedisAddr                        string        `mapstructure:"redisAddr"`
	DB                               *int          `mapstructure:"db"`
	Pass                             string        `mapstructure:"pass"`
	Concurrency                      *int          `mapstructure:"concurrency"`
	Queue                            string        `mapstructure:"queue"`
	Retention                        time.Duration `mapstructure:"retention"` // in hours
	RetryCount                       *int          `mapstructure:"retryCount"`
	Delay                            time.Duration `mapstructure:"delay"` // in seconds
	EmailInvitationTaskDelay         time.Duration `mapstructure:"emailInvitationTaskDelay"`
	EmailInvitationTaskRetryCount    *int          `mapstructure:"emailInvitationTaskRetryCount"`
	EmailInvitationTaskRetryDelay    time.Duration `mapstructure:"emailInvitationTaskRetryDelay"`
	EventReminderTaskRetryCount      *int          `mapstructure:"eventReminderTaskRetryCount"`
	EventReminderTaskRetryDelay      time.Duration `mapstructure:"eventReminderTaskRetryDelay"`
	EventReminderEmailTaskDelay      time.Duration `mapstructure:"eventReminderEmailTaskDelay"`
	EventReminderEmailTaskRetryCount *int          `mapstructure:"eventReminderEmailTaskRetryCount"`
	EventReminderEmailTaskRetryDelay time.Duration `mapstructure:"eventReminderEmailTaskRetryDelay"`
}

type RedisConfig struct {
	Host               string        `mapstructure:"host"`
	Port               *int          `mapstructure:"port"`
	Password           string        `mapstructure:"password"`
	DB                 *int          `mapstructure:"db"`
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
	Port     *int   `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type TelegramConfig struct {
	Token   string `mapstructure:"token"`
	Enabled *bool  `mapstructure:"enabled"`
}

type AWSConfig struct {
	Region          string        `mapstructure:"region"`
	AccessKeyID     string        `mapstructure:"access_key_id"`
	SecretAccessKey string        `mapstructure:"secret_access_key"`
	Endpoint        string        `mapstructure:"endpoint"`
	UseLocalStack   *bool         `mapstructure:"use_localstack"`
	ConfigService   ServiceConfig `mapstructure:"config_service"`
}

type KeycloakConfig struct {
	ServerURL    string `mapstructure:"server_url"`
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type ServiceConfig struct {
	Enabled *bool  `mapstructure:"enabled"`
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
	EnvAppDomain  = "APP_DOMAIN"

	EnvGeminiKey = "GEMINI_KEY"

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

	EnvTelegramToken   = "TELEGRAM_TOKEN"
	EnvTelegramEnabled = "TELEGRAM_ENABLED"

	EnvAPIKeyEncryptionSecret = "API_KEY_ENCRYPTION_SECRET"
)

// LoadConfig loads configuration from file, environment variables, or SSM Parameter Store
func LoadConfig() {
	defaultConfig := getDefaults()
	loadFromSsm(defaultConfig)
	loadFromEnv()
	printAllVipers()

	if err := viper.Unmarshal(&appConfig); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		os.Exit(1)
	}
}

func loadFromSsm(defaultConfig *Config) {
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
			useLocalStack = *defaultConfig.AWS.UseLocalStack
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
		var ssmConfigMap Config
		err = json.Unmarshal([]byte(*resp.Parameter.Value), &ssmConfigMap)
		if err != nil {
			panic("failed to unmarshal config from SSM: " + err.Error())
		}

		setViperFields(ssmConfigMap, "")
	}
}

func printAllVipers() {
	for _, s := range viper.AllKeys() {
		val := viper.Get(s)
		// if it's a pointer, print the value
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Ptr {
			if !v.IsNil() {
				fmt.Printf("%s=%v\n", s, v.Elem().Interface())
			} else {
				fmt.Printf("%s=<nil>\n", s)
			}
		} else {
			fmt.Printf("%s=%v\n", s, val)
		}
	}
}

func getDefaults() *Config {
	TruePointer := true
	FalsePointer := false
	conf := &Config{
		App: AppConfig{
			Name:       "NotificationManagement",
			Version:    "1.0.0",
			Port:       ToInt("8080"),
			Env:        "development",
			Encryption: "laeoGcA0ZFFsm3d9SUKevwG4VL4QN9Yi",
			Domain:     "https://github.com/tuhin47/NotificationManagement",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     ToInt("5432"),
			User:     "postgres",
			Password: "",
			Name:     "notification_management",
			SSLMode:  "disable",
		},
		AsynqConfig: AsynqConfig{
			RedisAddr: "",
			Pass:      "",
			Queue:     "",
		},
		Redis: RedisConfig{
			Host:              "localhost",
			Port:              ToInt("6379"),
			Password:          "",
			DB:                ToInt("0"),
			MandatoryPrefix:   "",
			AccessUuidPrefix:  "",
			RefreshUuidPrefix: "",
			UserPrefix:        "",
			PermissionPrefix:  "",
		},
		Email: EmailConfig{
			Url:      "",
			Timeout:  0,
			Host:     "localhost",
			Port:     ToInt("1025"),
			Username: "",
			Password: "",
			From:     "noreply@example.com",
		},
		AWS: AWSConfig{
			Region:          "us-east-1",
			AccessKeyID:     "test",
			SecretAccessKey: "test",
			Endpoint:        "http://localhost:4566",
			UseLocalStack:   &TruePointer,
			ConfigService: ServiceConfig{
				Enabled: &TruePointer,
				SSM:     "/myapp/config",
			},
		},
		Logger: LoggerConfig{
			Level:    "info",
			FilePath: "logs/app.log",
			Format:   "console",
		},
		Keycloak: KeycloakConfig{
			ServerURL:    "http://localhost:8081",
			Realm:        "gocloak",
			ClientID:     "gocloak",
			ClientSecret: "gocloak-secret",
		},
		Telegram: TelegramConfig{
			Token:   "",
			Enabled: &FalsePointer,
		},
	}
	setViperFields(conf, "")
	return conf
}

func setViperFields(conf interface{}, prefix string) {
	v := reflect.ValueOf(conf)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		key := field.Tag.Get("mapstructure")
		if key != "" {
			var curr = key
			if prefix != "" {
				curr = fmt.Sprintf("%s.%s", prefix, key)
			}
			if contains(properties, key) {
				setViperFields(value.Interface(), curr)
			} else {
				if value.Kind() == reflect.Ptr && value.IsNil() {
					continue
				}
				val := value.Interface()
				if val == nil {
					continue
				}
				if s, ok := val.(string); ok && s == "" {
					continue
				}
				if ps, ok := val.(*string); ok && (ps == nil || *ps == "") {
					continue
				}

				viper.Set(curr, val)
			}
		}
	}
}
func contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ToInt(s string) *int {
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}

func toBool(s string) *bool {
	if s == "" {
		return nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &b
}

func loadFromEnv() *Config {
	// App environment variables
	c := &Config{
		App: AppConfig{
			Name:       os.Getenv(EnvAppName),
			Version:    os.Getenv(EnvAppVersion),
			Port:       ToInt(os.Getenv(EnvAppPort)),
			Env:        os.Getenv(EnvAppEnv),
			Encryption: os.Getenv(EnvAPIKeyEncryptionSecret),
			Domain:     os.Getenv(EnvAppDomain),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv(EnvDBHost),
			Port:     ToInt(os.Getenv(EnvDBPort)),
			User:     os.Getenv(EnvDBUser),
			Password: os.Getenv(EnvDBPassword),
			Name:     os.Getenv(EnvDBName),
			SSLMode:  os.Getenv(EnvDBSSLMode),
		},
		AsynqConfig: AsynqConfig{},
		Redis: RedisConfig{
			Host:     os.Getenv(EnvRedisHost),
			Port:     ToInt(os.Getenv(EnvRedisPort)),
			Password: os.Getenv(EnvRedisPassword),
			DB:       ToInt(os.Getenv(EnvRedisDB)),
		},
		Email: EmailConfig{
			Host:     os.Getenv(EnvEmailHost),
			Port:     ToInt(os.Getenv(EnvEmailPort)),
			Username: os.Getenv(EnvEmailUsername),
			Password: os.Getenv(EnvEmailPassword),
			From:     os.Getenv(EnvEmailFrom),
		},
		AWS: AWSConfig{
			Region:          os.Getenv(EnvAWSRegion),
			AccessKeyID:     os.Getenv(EnvAWSAccessKeyID),
			SecretAccessKey: os.Getenv(EnvAWSSecretAccessKey),
			Endpoint:        os.Getenv(EnvAWSEndpoint),
			UseLocalStack:   toBool(os.Getenv(EnvAWSUseLocalStack)),
			ConfigService: ServiceConfig{
				Enabled: toBool(os.Getenv(EnvAWSConfigServiceEnabled)),
			},
		},
		Logger: LoggerConfig{
			Level:    os.Getenv(EnvLogLevel),
			FilePath: os.Getenv(EnvLogFilePath),
			Format:   os.Getenv(EnvLogFormat),
		},
		Keycloak: KeycloakConfig{
			ServerURL:    os.Getenv(EnvKeycloakServerURL),
			Realm:        os.Getenv(EnvKeycloakRealm),
			ClientID:     os.Getenv(EnvKeycloakClientID),
			ClientSecret: os.Getenv(EnvKeycloakClientSecret),
		},
		Telegram: TelegramConfig{
			Token:   os.Getenv(EnvTelegramToken),
			Enabled: toBool(os.Getenv(EnvTelegramEnabled)),
		},
		Development: DevelopmentConfig{
			GeminiKey: os.Getenv(EnvGeminiKey),
		},
	}
	setViperFields(c, "")
	return c
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

func Telegram() TelegramConfig {
	return appConfig.Telegram
}

// GetDSN returns the database connection string
func GetDSN() string {
	db := Database()
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, *db.Port, db.User, db.Password, db.Name, db.SSLMode)
}

// GetRedisAddr returns the Redis connection address
func GetRedisAddr() string {
	redis := Redis()
	return fmt.Sprintf("%s:%d", redis.Host, *redis.Port)
}

func Development() DevelopmentConfig {
	return appConfig.Development
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
	return *AWS().UseLocalStack
}
