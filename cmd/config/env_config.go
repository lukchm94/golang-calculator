package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type EnvVarsKeys string

const (
	PortEnvKey                  EnvVarsKeys = "PORT"
	AwsDefaultRegionEnvKey      EnvVarsKeys = "AWS_REGION"
	AwsAccessKeyEnvKey          EnvVarsKeys = "AWS_ACCESS_KEY_ID"
	AwsSecretKeyEnvKey          EnvVarsKeys = "AWS_SECRET_ACCESS_KEY"
	LocalstackEndpointUrlEnvKey EnvVarsKeys = "AWS_LOCALSTACK_ENDPOINT_URL"
	JwtSecretKeyEnvKey          EnvVarsKeys = "JWT_SECRET_KEY"
	JwtExpirationTimeEnvKey     EnvVarsKeys = "JWT_EXPIRATION_TIME"
	AppIssuerEnvKey             EnvVarsKeys = "APP_ISSUER"
	PostgresHostEnvKey          EnvVarsKeys = "POSTGRES_HOST"
	PostgresPortEnvKey          EnvVarsKeys = "POSTGRES_PORT"
	PostgresUserEnvKey          EnvVarsKeys = "POSTGRES_USER"
	PostgresPasswordEnvKey      EnvVarsKeys = "POSTGRES_PASSWORD"
	PostgresDbEnvKey            EnvVarsKeys = "POSTGRES_DB"
	Env                         EnvVarsKeys = "ENV"
)

type Configs struct {
	Port                  string
	AwsDefaultRegion      string
	AwsAccessKey          string
	AwsSecretKey          string
	LocalstackEndpointUrl string
	JwtSecretKey          string
	JwtExpirationTime     string
	AppIssuer             string
	PostgresHost          string
	PostgresPort          string
	PostgresUser          string
	PostgresPassword      string
	PostgresDb            string
	Env                   string
	AwsConfig             AwsConfig
}

func LoadConfigs(logger *slog.Logger) Configs {
	if err := godotenv.Load(); err != nil {
		logger.Error("No .env file found, using system environment")
	}

	env := getEnvVar(Env, string(DefaultEnv))

	if !ValidEnvironments(env).IsValid() {
		logger.Error("Invalid environment: %s", "env", env)
	}

	awsConfig := GetAwsConfig(ValidEnvironments(env), logger)

	return Configs{
		Port:                  getEnvVar(PortEnvKey, string(DefaultPort)),
		AwsDefaultRegion:      getEnvVar(AwsDefaultRegionEnvKey, string(EmptyString)),
		AwsAccessKey:          getEnvVar(AwsAccessKeyEnvKey, string(EmptyString)),
		AwsSecretKey:          getEnvVar(AwsSecretKeyEnvKey, string(EmptyString)),
		LocalstackEndpointUrl: getEnvVar(LocalstackEndpointUrlEnvKey, string(EmptyString)),
		JwtSecretKey:          getEnvVar(JwtSecretKeyEnvKey, string(EmptyString)),
		JwtExpirationTime:     getEnvVar(JwtExpirationTimeEnvKey, string(EmptyString)),
		AppIssuer:             getEnvVar(AppIssuerEnvKey, string(AppIssuer)),
		PostgresHost:          getEnvVar(PostgresHostEnvKey, string(EmptyString)),
		PostgresPort:          getEnvVar(PostgresPortEnvKey, string(EmptyString)),
		PostgresUser:          getEnvVar(PostgresUserEnvKey, string(EmptyString)),
		PostgresPassword:      getEnvVar(PostgresPasswordEnvKey, string(EmptyString)),
		PostgresDb:            getEnvVar(PostgresDbEnvKey, string(EmptyString)),
		Env:                   env,
		AwsConfig:             awsConfig,
	}
}

func getEnvVar(key EnvVarsKeys, defaultValue string) string {
	value := os.Getenv(string(key))
	if value == string(EmptyString) {
		return defaultValue
	}
	return value
}
