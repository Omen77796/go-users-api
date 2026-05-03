package config

import "os"

type Config struct {
	DBUrl         string
	RedisAddr     string
	RedisPassword string
	Port          string
}

func Load() *Config {
	return &Config{
		DBUrl:         os.Getenv("DATABASE_URL"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		Port:          os.Getenv("SERVER_PORT"),
	}
}
