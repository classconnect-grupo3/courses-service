package config

import "os"

type Config struct {
	DBUsername  string
	DBPassword  string
	DBName      string
	Host        string
	Port        string
	Environment string
}

func NewConfig() *Config {
	return &Config{
		DBUsername:  os.Getenv("DB_USERNAME"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		Host:        os.Getenv("HOST"),
		Port:        os.Getenv("PORT"),
		Environment: os.Getenv("ENVIRONMENT"),
	}
}
