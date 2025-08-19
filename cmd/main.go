package main

import (
	"log"

	"github.com/MapMinder/mapminder_backend/internal/config"
	infrastructure "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	"github.com/MapMinder/mapminder_backend/internal/server"
)

func main() {
	// 設定を読み込み
	cfg := config.Load()
	dbCfg := config.LoadDbConfig()
	dbConn := infrastructure.NewDBHandler(dbCfg)

	// サーバーを初期化
	srv := server.NewServer(cfg, dbConn)

	// サーバーを開始
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
