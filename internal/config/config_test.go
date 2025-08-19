package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	// setup
	t.Setenv("GO_ENV", "test")
	t.Setenv("PORT", "8080")

	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("expected port 8080 but got %s", cfg.Port)
	}
}

func TestLoadDBConfig(t *testing.T) {
	// setup
	t.Setenv("GO_ENV", "test")
	t.Setenv("DB_USER", "mapminder")
	t.Setenv("DB_PASSWORD", "mapminderpass")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_NAME", "mapminder")

	dbCfg := LoadDbConfig()

	if dbCfg.DBUser != "mapminder" {
		t.Errorf("expected port 8080 but got %s", dbCfg.DBUser)
	}

	if dbCfg.DBPassword != "mapminderpass" {
		t.Errorf("expected password mapminderpass but got %s", dbCfg.DBPassword)
	}

	if dbCfg.DBHost != "localhost" {
		t.Errorf("expected htt localhost but got %s", dbCfg.DBHost)
	}

	if dbCfg.DBPort != "3306" {
		t.Errorf("expected DBPort 3306 but got %s", dbCfg.DBPort)
	}

	if dbCfg.DBName != "mapminder" {
		t.Errorf("expected DBName mapminder but got %s", dbCfg.DBName)
	}
}
