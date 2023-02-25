package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Host string
	Port int
}

var AppConfig Config

// GetEnvStr get env string
func GetEnvStr(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// GetEnvInt get env int
func GetEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("invalid environment value env: %s val: %s", key, val)
	}

	return valInt
}

// initialize config
//
// To read from environment variable, os.Getenv("ENV_NAME")
func init() {
	AppConfig = Config{
		Host: GetEnvStr("APP_HOST", ""),
		Port: GetEnvInt("APP_PORT", 8000),
	}
}
