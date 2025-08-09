package config

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
