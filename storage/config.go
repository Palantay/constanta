package storage

import (
	"os"
)

type Config struct {
	DatabaseURI string
}

func NewConfig() *Config {
	return &Config{
		DatabaseURI: os.Getenv("database_uri"),
	}
}
