package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	Inbox    string
	Username string
	Password string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Host:     os.Getenv("IMAP_HOST_NAME"),
		Port:     os.Getenv("IMAP_PORT"),
		Inbox:    os.Getenv("EMAIL_INBOX"),
		Username: os.Getenv("IMAP_USERNAME"),
		Password: os.Getenv("IMAP_PASSWORD"),
	}

	if cfg.Host == "" || cfg.Port == "" || cfg.Inbox == "" || cfg.Username == "" || cfg.Password == "" {
		return nil, err
	}

	return cfg, nil
}
