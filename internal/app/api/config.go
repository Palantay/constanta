package api

import (
	"log"
	"os"

	"github.com/Palantay/constanta/storage"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	LoggerLevel string
	Storage     *storage.Config
}

func NewConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file. Error: %s", err)

	}

	return &Config{
		Port:        os.Getenv("port"),
		LoggerLevel: os.Getenv("logger_level"),
		Storage:     storage.NewConfig(),
	}

}
