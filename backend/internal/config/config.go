package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBURL     string
	REDIS_URL string
	JWTSecret string
	MasterKey string
	EthRPCURL string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		DBURL:     getEnv("DATABASE_URL", ""),
		REDIS_URL: getEnv("REDIS_URL", "localhost:6379"),
		JWTSecret: getEnv("JWT_SECRET", ""),
		MasterKey: getEnv("MASTER_KEY", ""),
		EthRPCURL: getEnv("ETH_RPC_URL", "http://localhost:8545"),
	}

	// Validate required environment variables for production
	if cfg.DBURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	if cfg.MasterKey == "" {
		log.Fatal("MASTER_KEY environment variable is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
