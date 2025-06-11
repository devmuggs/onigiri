package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DBConnectionString string
}

func New(ctx context.Context, logger *zap.Logger) (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found, relying on environment variables", zap.Error(err))
	}

	dbConnectionString := os.Getenv("DATABASE_URL")
	if dbConnectionString == "" {
		return nil, fmt.Errorf("unable to find env 'DATABASE_URL'")
	}

	return &Config{DBConnectionString: dbConnectionString}, nil
}
