package server

import (
	"meme-generator/internal/config"
	"meme-generator/internal/meme"
	"meme-generator/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine

	memeManager *meme.MemeManager
	cfg         *config.Config
	store       storage.Store
}

func NewServer(secretKey string, memeManager *meme.MemeManager, cfg *config.Config, store storage.Store) *Server {
	// CORS
	ginConfig := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}

	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20

	sessionStore := cookie.NewStore([]byte(secretKey))
	sessionStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24, // 1 day
		Secure:   false,
		HttpOnly: true,
	})

	router.Use(
		cors.New(ginConfig),
		sessions.Sessions("session", sessionStore),
	)

	server := &Server{
		router:      router,
		memeManager: memeManager,
		cfg:         cfg,
		store:       store,
	}

	server.registerRoutes()

	return server
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
