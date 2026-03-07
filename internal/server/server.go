package server

import (
	"fmt"

	"github.com/MapMinder/mapminder_backend/internal/config"
	router "github.com/MapMinder/mapminder_backend/internal/httpServer"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	engine *gin.Engine
	config *config.Config
}

func NewServer(cfg *config.Config, dbHandler *gorm.DB) *Server {
	// set environment
	SetServerEnvironment(cfg)

	/*
			  NOTE: gin.Default() は engine.Use(Logger(), Recovery())両方初期化してしまう
		    このアップリではzapLoggerを使用しているためgin.Default()ではなくgin.New()で大丈夫
	*/
	engine := gin.New()
	engine.Use(gin.Recovery())
	printBanner(cfg.Port, cfg.Environment)

	// add request logger
	engine.Use(logger.MiddlewareLogger(logger.Get()))

	// ルートを設定
	router.NewRoutes(engine, dbHandler)

	return &Server{
		engine: engine,
		config: cfg,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Port)

	return s.engine.Run(addr)
}

// SetServerEnvironment
func SetServerEnvironment(cfg *config.Config) {
	switch cfg.Environment {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "development":
		gin.SetMode(gin.DebugMode)
	}
}

func printBanner(port string, environment string) {
	banner := `
 _______  _______  _______  _______ _________ _        ______   _______  _______ 
(       )(  ___  )(  ____ )(       )\__   __/( (    /|(  __  \ (  ____ \(  ____ )
| () () || (   ) || (    )|| () () |   ) (   |  \  ( || (  \  )| (    \/| (    )|
| || || || (___) || (____)|| || || |   | |   |   \ | || |   ) || (__    | (____)|
| |(_)| ||  ___  ||  _____)| |(_)| |   | |   | (\ \) || |   | ||  __)   |     __)
| |   | || (   ) || (      | |   | |   | |   | | \   || |   ) || (      | (\ (   
| )   ( || )   ( || )      | )   ( |___) (___| )  \  || (__/  )| (____/\| ) \ \__
|/     \||/     \||/       |/     \|\_______/|/    )_)(______/ (_______/|/   \__/
                                                                                 

`
	fmt.Println(banner)
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Server running at http://localhost:%s\n", port)
	fmt.Println()
}
