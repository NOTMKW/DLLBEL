package config

import (
	"os"
)

type Config struct {
	Port 	string
	RedisAddr 	string
	RedisPass 	string
	RedisDB 	int
	EventBuffer int
	Workers 	int
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass: getEnv("REDIS_PASS", ""),
		RedisDB: 0,
		EventBuffer: 10000,
		Workers: 10,	
	}
}
	func getEnv(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
	return defaultValue
	}