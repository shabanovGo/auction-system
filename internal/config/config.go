package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	HTTP     HTTPConfig
	GRPC     GRPCConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type HTTPConfig struct {
	Host string
	Port int
}

type GRPCConfig struct {
	Host string
	Port int
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}
	
	cfg.Server.Host = os.Getenv("SERVER_HOST")
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}
	cfg.Server.Port = port
	
	readTimeout, err := time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_READ_TIMEOUT: %w", err)
	}
	cfg.Server.ReadTimeout = readTimeout
	
	writeTimeout, err := time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %w", err)
	}
	cfg.Server.WriteTimeout = writeTimeout

	cfg.Database.Host = os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	cfg.Database.Port = dbPort
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.DBName = os.Getenv("DB_NAME")
	cfg.Database.SSLMode = os.Getenv("DB_SSLMODE")

	cfg.HTTP.Host = getEnvOrDefault("HTTP_HOST", "0.0.0.0")
	httpPort, err := strconv.Atoi(getEnvOrDefault("HTTP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
	}
	cfg.HTTP.Port = httpPort

	cfg.GRPC.Host = getEnvOrDefault("GRPC_HOST", "0.0.0.0")
	grpcPort, err := strconv.Atoi(getEnvOrDefault("GRPC_PORT", "50051"))
	if err != nil {
		return nil, fmt.Errorf("invalid GRPC_PORT: %w", err)
	}
	cfg.GRPC.Port = grpcPort

	return cfg, nil
}

func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
