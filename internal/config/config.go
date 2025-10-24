package config

import (
	"os"
	"sync"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port        string
	Logs        LogConfig
	Env         string
	PostgresURI string
}

var (
	cfg  *Config
	once sync.Once
)

type LogConfig struct {
	Level    string
	SaveLogs bool
}

func LoadConfig() {
	// runs exactly once
	once.Do(func() {
		cfg = &Config{
			Port: ":" + os.Getenv("PORT"),
			Logs: LogConfig{
				Level:    getEnvWithDefault("LOG_LEVEL", "ERROR"),
				SaveLogs: false, // set false for now
			},
			PostgresURI: os.Getenv("DATABASE_URL"),
		}
	})
}
func Get() *Config {
	if cfg == nil {
		panic("Config not loaded, call config.Load()") // may not need to panic here
	}
	return cfg
}

// will attempt to get a .env value if non is found returns a default value
func getEnvWithDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val

}
