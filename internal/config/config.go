package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

func Load() (Config, error) {

	_ = godotenv.Load()

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return Config{}, fmt.Errorf("REDIS_DB must be a valid integer")
	}

	cfg := Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8000"),

		MySQLHost:     getEnv("MYSQL_HOST", ""),
		MySQLPort:     getEnv("MYSQL_PORT", ""),
		MySQLUser:     getEnv("MYSQL_USER", ""),
		MySQLPassword: getEnv("MYSQL_PASSWORD", ""),
		MySQLDatabase: getEnv("MYSQL_DATABASE", ""),

		RedisHost:     getEnv("REDIS_HOST", ""),
		RedisPort:     getEnv("REDIS_PORT", ""),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.AppEnv == "" {
		return fmt.Errorf("APP_ENV is required")
	}

	if c.AppPort == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	if c.MySQLHost == "" {
		return fmt.Errorf("MYSQL_HOST is required")
	}

	if c.MySQLPort == "" {
		return fmt.Errorf("MYSQL_PORT is required")
	}

	if c.MySQLUser == "" {
		return fmt.Errorf("MYSQL_USER is required")
	}

	if c.MySQLPassword == "" {
		return fmt.Errorf("MYSQL_PASSWORD is required")
	}

	if c.MySQLDatabase == "" {
		return fmt.Errorf("MYSQL_DATABASE is required")
	}

	if c.RedisHost == "" {
		return fmt.Errorf("REDIS_HOST is required")
	}

	if c.RedisPort == "" {
		return fmt.Errorf("REDIS_PORT is required")
	}

	if c.RedisDB < 0 {
		return fmt.Errorf("REDIS_DB cannot be negative")
	}

	return nil
}

func (c Config) MySQLDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		c.MySQLUser,
		c.MySQLPassword,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDatabase,
	)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
