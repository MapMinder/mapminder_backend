package config

import (
	"log"
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
		log.Fatal("failed to load .env file")
	}
	envLoaded = true
}

// 設定ファイルの読み込み
func Load() *Config {
	// setup
	LoadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
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

	// gets db user if provided kills process with log.fatal
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER environment variable is not set")
	}

	// gets db password if provided kills process with log.fatal
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is not set")
	}

	// gets db host if provided kills process with log.fatal
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST environment variable is not set")
	}

	// gets db port if provided kills process with log.fatal
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT environment variable is not set")
	}

	// gets db name if provided kills process with log.fatal
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is not set")
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
