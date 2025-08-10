package config

import (
	"NotificationManagement/config/helper"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
)

type Config struct {
	App         AppConfig         `mapstructure:"app" tag:"obj"`
	Database    DatabaseConfig    `mapstructure:"database" tag:"obj"`
	Asynq       AsynqConfig       `mapstructure:"asynq" tag:"obj"`
	Redis       RedisConfig       `mapstructure:"redis" tag:"obj"`
	Email       EmailConfig       `mapstructure:"email" tag:"obj"`
	AWS         AWSConfig         `mapstructure:"aws" tag:"obj"`
	Logger      LoggerConfig      `mapstructure:"logger" tag:"obj"`
	Keycloak    KeycloakConfig    `mapstructure:"keycloak" tag:"obj"`
	Telegram    TelegramConfig    `mapstructure:"telegram" tag:"obj"`
	Development DevelopmentConfig `mapstructure:"development" tag:"obj"`
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
	SSLMode  string `mapstructure:"sslMode"`
}

type AsynqConfig struct {
	RedisAddr                   string `mapstructure:"redisaddr"`
	DB                          *int   `mapstructure:"db"`
	Pass                        string `mapstructure:"pass"`
	Concurrency                 *int   `mapstructure:"concurrency"`
	Queue                       string `mapstructure:"queue"`
	Retention                   *int   `mapstructure:"retention"` // in Hours
	RetryCount                  *int   `mapstructure:"retryCount"`
	EventReminderTaskRetryCount *int   `mapstructure:"eventReminderTaskRetryCount"`
	EventReminderTaskRetryDelay *int   `mapstructure:"eventReminderTaskRetryDelay"` // in seconds
}

type RedisConfig struct {
	Host               string `mapstructure:"host"`
	Port               *int   `mapstructure:"port"`
	Password           string `mapstructure:"password"`
	DB                 *int   `mapstructure:"db"`
	MandatoryPrefix    string `mapstructure:"mandatoryPrefix"`
	AccessUuidPrefix   string `mapstructure:"accessUuidPrefix"`
	RefreshUuidPrefix  string `mapstructure:"refreshUuidPrefix"`
	UserPrefix         string `mapstructure:"userPrefix"`
	PermissionPrefix   string `mapstructure:"permissionPrefix"`
	UserCacheTTL       *int   `mapstructure:"userCacheTTL"`
	PermissionCacheTTL *int   `mapstructure:"permissionCacheTTL"`
}

type EmailConfig struct {
	Url      string
	Timeout  *int
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
	AccessKeyID     string        `mapstructure:"accessKeyID"`
	SecretAccessKey string        `mapstructure:"secretAccessKey"`
	Endpoint        string        `mapstructure:"endpoint"`
	UseLocalStack   *bool         `mapstructure:"useLocalstack"`
	ConfigService   ServiceConfig `mapstructure:"configService" tag:"obj"`
}

type KeycloakConfig struct {
	ServerURL    string `mapstructure:"serverUrl"`
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"clientId"`
	ClientSecret string `mapstructure:"clientSecret"`
}

type ServiceConfig struct {
	Enabled *bool  `mapstructure:"enabled"`
	SSM     string `mapstructure:"ssm"`
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"filePath"`
	Format   string `mapstructure:"format"`
}

var (
	appConfig *Config
)

func LoadConfig() {
	defaultConfig := loadDefaults()
	cnf := &defaultConfig.AWS
	cnf = cnf.MergeWithEnv()
	client, err := NewAWSClient(cnf)
	if err == nil {
		client.loadFromSsm()
	} else {
		fmt.Println("AWS Client error", err.Error())
	}

	loadFromEnv()

	if err := viper.Unmarshal(&appConfig); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		os.Exit(1)
	}
	printAllVipers()
}

func printAllVipers() {
	if IsDevelopment() {
		return
	}
	for _, s := range viper.AllKeys() {
		val := viper.Get(s)
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

func loadDefaults() *Config {
	TruePointer := true
	FalsePointer := false
	conf := &Config{
		App: AppConfig{
			Name:       "NotificationManagement",
			Version:    "1.0.0",
			Port:       helper.ToInt("8080"),
			Env:        "development",
			Encryption: "laeoGcA0ZFFsm3d9SUKevwG4VL4QN9Yi",
			Domain:     "https://github.com/tuhin47/NotificationManagement",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     helper.ToInt("54322"),
			User:     "user",
			Password: "password",
			Name:     "notification_management",
			SSLMode:  "disable",
		},
		Asynq: AsynqConfig{
			RedisAddr:                   "127.0.0.1:6379",
			DB:                          helper.ToInt("15"),
			Pass:                        "*****",
			Concurrency:                 helper.ToInt("10"),
			Queue:                       "notification-management",
			Retention:                   helper.ToInt("168"),
			RetryCount:                  helper.ToInt("25"),
			EventReminderTaskRetryCount: helper.ToInt("5"),
			EventReminderTaskRetryDelay: helper.ToInt("30"),
		},
		Redis: RedisConfig{
			Host:               "127.0.0.1",
			Port:               helper.ToInt("6379"),
			Password:           "",
			DB:                 helper.ToInt("0"),
			MandatoryPrefix:    "n_m_",
			AccessUuidPrefix:   "a-uuid_",
			RefreshUuidPrefix:  "r-uuid_",
			UserPrefix:         "user_",
			PermissionPrefix:   "permissions_",
			UserCacheTTL:       helper.ToInt("3600"),
			PermissionCacheTTL: helper.ToInt("86400"),
		},
		Email: EmailConfig{
			Url:      "",
			Timeout:  helper.ToInt("0"),
			Host:     "localhost",
			Port:     helper.ToInt("1025"),
			Username: "Admin",
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
			Format:   "",
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
		Development: DevelopmentConfig{
			GeminiKey: "",
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
		obj := field.Tag.Get("tag")
		if key != "" {
			var curr = key
			if prefix != "" {
				curr = fmt.Sprintf("%s.%s", prefix, key)
			}
			if obj == "obj" {
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

func loadFromEnv() *Config {
	c := &Config{
		App: AppConfig{
			Name:       os.Getenv(EnvAppName),
			Version:    os.Getenv(EnvAppVersion),
			Port:       helper.ToInt(os.Getenv(EnvAppPort)),
			Env:        os.Getenv(EnvAppEnv),
			Encryption: os.Getenv(EnvAPIKeyEncryptionSecret),
			Domain:     os.Getenv(EnvAppDomain),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv(EnvDBHost),
			Port:     helper.ToInt(os.Getenv(EnvDBPort)),
			User:     os.Getenv(EnvDBUser),
			Password: os.Getenv(EnvDBPassword),
			Name:     os.Getenv(EnvDBName),
			SSLMode:  os.Getenv(EnvDBSSLMode),
		},
		Asynq: AsynqConfig{},
		Redis: RedisConfig{
			Host:     os.Getenv(EnvRedisHost),
			Port:     helper.ToInt(os.Getenv(EnvRedisPort)),
			Password: os.Getenv(EnvRedisPassword),
			DB:       helper.ToInt(os.Getenv(EnvRedisDB)),
		},
		Email: EmailConfig{
			Host:     os.Getenv(EnvEmailHost),
			Port:     helper.ToInt(os.Getenv(EnvEmailPort)),
			Username: os.Getenv(EnvEmailUsername),
			Password: os.Getenv(EnvEmailPassword),
			From:     os.Getenv(EnvEmailFrom),
		},
		AWS: AWSConfig{
			Region:          os.Getenv(EnvAWSRegion),
			AccessKeyID:     os.Getenv(EnvAWSAccessKeyID),
			SecretAccessKey: os.Getenv(EnvAWSSecretAccessKey),
			Endpoint:        os.Getenv(EnvAWSEndpoint),
			UseLocalStack:   helper.ToBool(os.Getenv(EnvAWSUseLocalStack)),
			ConfigService: ServiceConfig{
				Enabled: helper.ToBool(os.Getenv(EnvAWSConfigServiceEnabled)),
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
			Enabled: helper.ToBool(os.Getenv(EnvTelegramEnabled)),
		},
		Development: DevelopmentConfig{
			GeminiKey: os.Getenv(EnvGeminiKey),
		},
	}
	setViperFields(c, "")
	return c
}

func (c *AWSConfig) MergeWithEnv() *AWSConfig {
	return &AWSConfig{
		Region:          *helper.FirstNonEmpty(os.Getenv(EnvAWSRegion), c.Region),
		AccessKeyID:     *helper.FirstNonEmpty(os.Getenv(EnvAWSAccessKeyID), c.AccessKeyID),
		SecretAccessKey: *helper.FirstNonEmpty(os.Getenv(EnvAWSSecretAccessKey), c.SecretAccessKey),
		Endpoint:        *helper.FirstNonEmpty(os.Getenv(EnvAWSEndpoint), c.Endpoint),
		UseLocalStack:   *helper.FirstNonEmpty(helper.ToBool(os.Getenv(EnvAWSUseLocalStack)), c.UseLocalStack),
		ConfigService: ServiceConfig{
			SSM:     *helper.FirstNonEmpty(os.Getenv(EnvConfigSSMParam), c.ConfigService.SSM),
			Enabled: *helper.FirstNonEmpty(helper.ToBool(os.Getenv(EnvAWSConfigServiceEnabled)), c.ConfigService.Enabled),
		},
	}
}

func App() AppConfig {
	return appConfig.App
}

func Database() DatabaseConfig {
	return appConfig.Database
}

func Asynq() AsynqConfig {
	return appConfig.Asynq
}

func Redis() RedisConfig {
	return appConfig.Redis
}

func Email() EmailConfig {
	return appConfig.Email
}

func AWS() AWSConfig {
	return appConfig.AWS
}

func Logger() LoggerConfig {
	return appConfig.Logger
}

func Keycloak() KeycloakConfig {
	return appConfig.Keycloak
}

func Telegram() TelegramConfig {
	return appConfig.Telegram
}

func GetDSN() string {
	db := Database()
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, *db.Port, db.User, db.Password, db.Name, db.SSLMode)
}

func GetRedisAddr() string {
	redis := Redis()
	return fmt.Sprintf("%s:%d", redis.Host, *redis.Port)
}

func Development() DevelopmentConfig {
	return appConfig.Development
}

func IsDevelopment() bool {
	return App().Env == "development"
}

func IsProduction() bool {
	return App().Env == "production"
}
