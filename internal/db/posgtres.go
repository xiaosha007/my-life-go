package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var DB *sqlx.DB

func InitConnection(config *Config) error {
	dataSourceName := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s", config.User, config.DBName, config.SSLMode, config.Password, config.Host)

	var err error

	DB, err = sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}

	// Set database connection pool settings
	DB.SetMaxOpenConns(25)                 // Maximum number of open connections
	DB.SetMaxIdleConns(25)                 // Maximum number of idle connections
	DB.SetConnMaxLifetime(5 * time.Minute) // Lifetime of a connection

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to DB!")
	}

	return nil
}

func CloseConnection() error {

	if DB != nil {
		err := DB.Close()
		if err != nil {
			return fmt.Errorf("failed to close the database: %w", err)
		}
		log.Println("DB closed gracefully...")
	}
	return nil
}
