package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  int
	JWTSecret   string
	JWTExpiry   time.Duration
	Env         string
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	FrontendURL string
}

var AppConfig *Config

func Load() error {
	godotenv.Load()

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))

	AppConfig = &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      dbPort,
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "family_tree"),
		ServerPort:  serverPort,
		JWTSecret:   getEnv("JWT_SECRET", "secret"),
		JWTExpiry:   jwtExpiry,
		Env:         getEnv("ENV", "development"),
		SMTPHost:    getEnv("SMTP_HOST", ""),
		SMTPPort:    smtpPort,
		SMTPUser:    getEnv("SMTP_USER", ""),
		SMTPPass:    getEnvWithFallback("SMTP_PASS", "SMTP_PASSWORD", ""),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	return nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvWithFallback(primaryKey, fallbackKey, defaultVal string) string {
	if value, exists := os.LookupEnv(primaryKey); exists {
		return value
	}
	if value, exists := os.LookupEnv(fallbackKey); exists {
		return value
	}
	return defaultVal
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}
