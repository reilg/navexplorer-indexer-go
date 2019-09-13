package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Storage  StorageConfig
	Navcoind NavcoindConfig
}

type StorageConfig struct {
	Path string
}

type NavcoindConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}


func Init() {
	err := godotenv.Load()
	if err != nil {
 		log.Fatal(err)
	}
}

func Get() *Config {
	return &Config{
		Storage: StorageConfig{
			Path: getString("STORAGE_PATH", "/data"),
		},
		Navcoind: NavcoindConfig{
			Host:     getString("NAVCOIND_HOST", ""),
			Port:     getInt("NAVCOIND_PORT", 8332),
			User:     getString("NAVCOIND_USER", "user"),
			Password: getString("NAVCOIND_PASSWORD", "password"),
		},
	}
}

func getString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getInt(key string, defaultValue int) int {
	valStr := getString(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}

	return defaultValue
}

func getBool(key string, defaultValue bool) bool {
	valStr := getString(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultValue
}

func getSlice(key string, defaultVal []string, sep string) []string {
	valStr := getString(key, "")
	if valStr == "" {
		return defaultVal
	}

	return strings.Split(valStr, sep)
}