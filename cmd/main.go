package main

import (
	"log"

	"github.com/MapMinder/mapminder_backend/internal/config"
	"github.com/MapMinder/mapminder_backend/internal/server"
)

func main() {
	// 設定を読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// サーバーを初期化
	srv := server.NewServer(cfg)

	// サーバーを開始
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
