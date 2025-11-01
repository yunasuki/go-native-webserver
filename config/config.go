package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

type ServerConfig struct {
	AllowedOrigins      []string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	QueueWorkMaxCount   int

	Database ServerDatabaseConfig
}

type ServerDatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

var (
	configInstance *ServerConfig
	configMutex    sync.RWMutex
)

// GetServerConfig returns the singleton instance of ServerConfig.
func GetServerConfig() *ServerConfig {
	configMutex.RLock()
	if configInstance != nil {
		defer configMutex.RUnlock()
		return configInstance
	}
	configMutex.RUnlock()
	configMutex.Lock()
	defer configMutex.Unlock()
	if configInstance == nil {
		configInstance = loadConfig()
	}
	return configInstance
}

// ReloadServerConfig reloads the configuration from environment variables.
func ReloadServerConfig() {
	configMutex.Lock()
	defer configMutex.Unlock()
	configInstance = loadConfig()
}

func loadConfig() *ServerConfig {
	allowedOriginsRaw := getEnv("ALLOWED_ORIGINS", "*")
	allowedOrigins := strings.Split(allowedOriginsRaw, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}
	readTimeout := getEnvAsInt("READ_TIMEOUT_SECONDS", 10)
	writeTimeout := getEnvAsInt("WRITE_TIMEOUT_SECONDS", 30)
	queueWorkMaxCount := getEnvAsInt("QUEUE_WORK_MAX_COUNT", 100)

	databaseConfig := ServerDatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		Name:     getEnv("DB_NAME", "appdb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	return &ServerConfig{
		AllowedOrigins:      allowedOrigins,
		ReadTimeoutSeconds:  readTimeout,
		WriteTimeoutSeconds: writeTimeout,
		QueueWorkMaxCount:   queueWorkMaxCount,
		Database:            databaseConfig,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultVal int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultVal
}
