package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/xiaosha007/my-life-go/internal/db"
	"github.com/xiaosha007/my-life-go/internal/monitoring"
	"github.com/xiaosha007/my-life-go/pkg/transformer"
)

type Config struct {
	DB        *db.Config
	Statsd    *monitoring.StatsdClientConfig
	JwtSecret string
}

func GetConfig() *Config {
	config := &Config{
		DB: &db.Config{Host: os.Getenv("DB_HOST"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Statsd: &monitoring.StatsdClientConfig{
			Namespace: os.Getenv("STATSD_NAMESPACE"),
			Host:      os.Getenv("STATSD_HOST"),
			Port:      getEnvAsInt("STATSD_PORT", 8125),
		},
		JwtSecret: os.Getenv("JWT_SECRET"),
	}

	return config
}

func Init() error {
	// Construct the path to the .env file at the root of the project
	rootPath, err := filepath.Abs("../../.env")
	if err != nil {
		return fmt.Errorf("error resolving .env file path: %s", err)
	}

	err = godotenv.Load(rootPath)

	if err != nil {
		return err
	}

	log.Println("Config loaded!")

	return nil
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultVal
	}
	intEnv, err := transformer.StringToInt(valueStr)
	if err != nil {
		log.Printf("Invalid integer for %s: %v. Using default value: %d", name, err, defaultVal)
		return defaultVal
	}
	return intEnv
}
