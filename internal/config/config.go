// Package config handles application configuration by loading environment variables from a .env file.
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var configuration *Config

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type Config struct {
	Version      string
	ServiceName  string
	Port         int64
	JWTSecretKey string
	DB           *DBConfig
}

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Must have an env file", err)
		os.Exit(1)
	}

	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("VERSION is required.")
		os.Exit(1)
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		fmt.Println("SERVICE_NAME is required")
		os.Exit(1)
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		fmt.Println("PORT is required")
		os.Exit(1)
	}

	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		fmt.Println("PORT must be a number.")
		os.Exit(1)
	}

	jwtSecreKey := os.Getenv("JWT_SECRET_KEY")

	if jwtSecreKey == "" {
		fmt.Println("JWT_SECRET_KEY is required.")
		os.Exit(1)
	}

	// DB

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		fmt.Println("DB_HOST is required.")
		os.Exit(1)
	}

	dbPORTStr := os.Getenv("DB_PORT")
	if dbPORTStr == "" {
		fmt.Println("DB_PORT is required")
		os.Exit(1)
	}

	dbPORT, err := strconv.ParseInt(dbPORTStr, 10, 64)
	if err != nil {
		fmt.Println("DB_PORT must be a number.")
		os.Exit(1)
	}

	dbName := os.Getenv("DB_NAME")

	if dbName == "" {
		fmt.Println("DB_NAME is required.")
		os.Exit(1)
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		fmt.Println("DB_USER is required.")
		os.Exit(1)
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		fmt.Println("DB_PASSWORD is required.")
		os.Exit(1)
	}

	configuration = &Config{
		Version:      version,
		ServiceName:  serviceName,
		Port:         port,
		JWTSecretKey: jwtSecreKey,
		DB: &DBConfig{
			Host:     dbHost,
			Port:     int(dbPORT),
			Name:     dbName,
			User:     dbUser,
			Password: dbPass,
		},
	}
}

func GetConfig() *Config {
	if configuration == nil {
		loadConfig()
	}
	return configuration
}
