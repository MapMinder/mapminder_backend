package config

import (
	"os"

	"github.com/joho/godotenv"
)

var envLoaded = false

type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

type Config struct {
	Port string
}

// load environment variables
func LoadEnv() {
	if envLoaded {
		return
	}

	// test 環境の場合は必ず環境変数を初期化せずreturnする
	if os.Getenv("GO_ENV") == "test" {
		envLoaded = true
		return
	}

	if err := godotenv.Load(); err != nil {
		panic("failed to load .env file: %w")
	}
	envLoaded = true
}

// 設定ファイルの読み込み
func Load() *Config {
	// setup
	LoadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}

	config := &Config{
		Port: port,
	}

	return config
}

// setup db config
func LoadDbConfig() *DBConfig {
	// setup
	LoadEnv()

	// gets db user if provided else panics
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		panic("DB_USER environment variable is not set")
	}

	// gets db password if provided else panics
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		panic("DB_PASSWORD environment variable is not set")
	}

	// gets db host if provided else panics
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		panic("DB_HOST environment variable is not set")
	}

	// gets db port if provided else panics
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		panic("DB_PORT environment variable is not set")
	}

	// gets db name if provided else panics
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		panic("DB_NAME environment variable is not set")
	}

	dbConfig := &DBConfig{
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBName:     dbName,
	}

	return dbConfig
}
