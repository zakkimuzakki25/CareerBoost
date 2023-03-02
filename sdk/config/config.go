package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Interface interface {
	Load(envPath string) error
	Get(key string) string
}

type config struct{}

func Init() Interface {
	return &config{}
}

func (c *config) Load(envPath string) error {
	return godotenv.Load(envPath)
}

func (c *config) Get(key string) string {
	return os.Getenv(key)
}
