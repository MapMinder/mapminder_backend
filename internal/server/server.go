package server

import (
	"fmt"
	"log"

	"github.com/MapMinder/mapminder_backend/internal/config"
	router "github.com/MapMinder/mapminder_backend/internal/httpServer"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	engine := gin.Default()

	// ルートを設定
	router.NewRoutes(engine)

	return &Server{
		engine: engine,
		config: cfg,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.Printf("Server starting on port %s", s.config.Port)

	return s.engine.Run(addr)
}
