package config

type ServerDefaultConfig string

const (
	DefaultPort ServerDefaultConfig = "8080"
	DefaultEnv  ServerDefaultConfig = "dev"
)

type Utils string

const (
	EmptyString Utils = ""
)
