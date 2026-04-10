package config

import (
	"log"
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
	PostgresHostEnvKey          EnvVarsKeys = "POSTGRES_HOST"
	PostgresPortEnvKey          EnvVarsKeys = "POSTGRES_PORT"
	PostgresUserEnvKey          EnvVarsKeys = "POSTGRES_USER"
	PostgresPasswordEnvKey      EnvVarsKeys = "POSTGRES_PASSWORD"
	PostgresDbEnvKey            EnvVarsKeys = "POSTGRES_DB"
)

type Configs struct {
	Port                  string
	AwsDefaultRegion      string
	AwsAccessKey          string
	AwsSecretKey          string
	LocalstackEndpointUrl string
	PostgresHost          string
	PostgresPort          string
	PostgresUser          string
	PostgresPassword      string
	PostgresDb            string
}

func LoadConfigs() Configs {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}
	return Configs{
		Port:                  getEnvVar(PortEnvKey, string(DefaultPort)),
		AwsDefaultRegion:      getEnvVar(AwsDefaultRegionEnvKey, string(EmptyString)),
		AwsAccessKey:          getEnvVar(AwsAccessKeyEnvKey, string(EmptyString)),
		AwsSecretKey:          getEnvVar(AwsSecretKeyEnvKey, string(EmptyString)),
		LocalstackEndpointUrl: getEnvVar(LocalstackEndpointUrlEnvKey, string(EmptyString)),
		PostgresHost:          getEnvVar(PostgresHostEnvKey, string(EmptyString)),
		PostgresPort:          getEnvVar(PostgresPortEnvKey, string(EmptyString)),
		PostgresUser:          getEnvVar(PostgresUserEnvKey, string(EmptyString)),
		PostgresPassword:      getEnvVar(PostgresPasswordEnvKey, string(EmptyString)),
		PostgresDb:            getEnvVar(PostgresDbEnvKey, string(EmptyString)),
	}
}

func getEnvVar(key EnvVarsKeys, defaultValue string) string {
	value := os.Getenv(string(key))
	if value == string(EmptyString) {
		return defaultValue
	}
	return value
}
