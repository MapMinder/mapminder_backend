package main

import (
	"github.com/MapMinder/mapminder_backend/internal/config"
	infrastructure "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/MapMinder/mapminder_backend/internal/server"
	"github.com/MapMinder/mapminder_backend/shared/validator"
)

func main() {
	// 設定を読み込み
	cfg := config.Load()

	// logger を初期化
	zapCfg := config.LoadZapConfig()
	logger.Init(zapCfg)

	dbCfg := config.LoadDbConfig()
	dbConn := infrastructure.NewDBHandler(dbCfg)

	// サーバーを初期化
	srv := server.NewServer(cfg, dbConn)

	// validator
	validator.Init()

	// サーバーを開始
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
