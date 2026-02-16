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
)

type Configs struct {
	Port                  string
	AwsDefaultRegion      string
	AwsAccessKey          string
	AwsSecretKey          string
	LocalstackEndpointUrl string
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
	}
}

func getEnvVar(key EnvVarsKeys, defaultValue string) string {
	value := os.Getenv(string(key))
	if value == string(EmptyString) {
		return defaultValue
	}
	return value
}
