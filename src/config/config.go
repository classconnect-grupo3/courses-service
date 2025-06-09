package config

import "os"

type Config struct {
	DBUsername   string
	DBPassword   string
	DBName       string
	DBURI        string
	Host         string
	Port         string
	Environment  string
	GeminiApiKey string
}

func NewConfig() *Config {
	return &Config{
		DBUsername:   os.Getenv("DB_USERNAME"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		DBURI:        os.Getenv("DB_URI"),
		Host:         os.Getenv("HOST"),
		Port:         os.Getenv("PORT"),
		Environment:  os.Getenv("ENVIRONMENT"),
		GeminiApiKey: os.Getenv("GEMINI_API_KEY"),
	}
}
