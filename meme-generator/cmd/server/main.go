package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"meme-generator/internal/config"
	"meme-generator/internal/meme"
	"meme-generator/internal/server"
	"meme-generator/internal/storage"
	"meme-generator/log"

	"go.uber.org/zap"
	"moul.io/zapgorm2"
)

func main() {
	if err := os.Chdir(filepath.Dir(appFilePath())); err != nil {
		mustInit(fmt.Errorf("failed to set working directory: %w", err))
	}

	logger := log.InitLogger(config.DebugMode())
	zapGormLogger := zapgorm2.New(logger)

	configManager := config.NewConfigManager()

	if err := configManager.Load("."); err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	cfg := configManager.Config

	db, err := storage.NewDatabase(cfg.Database, zapGormLogger)
	if err != nil {
		log.Fatal("failed to initialize database", zap.Error(err))
	}

	store := storage.NewStorage(db)

	memeManager, err := meme.NewMemeManager(cfg.App.DataDir, store)
	if err != nil {
		log.Fatal("failed to initialize meme manager", zap.Error(err))
	}

	httpAddr := net.JoinHostPort("", strconv.Itoa(cfg.Web.Port))

	httpServer := server.NewServer(httpAddr, memeManager, cfg, store)

	log.Fatal("failed to start server", zap.Error(httpServer.Run(httpAddr)))
}

func mustInit(err error) {
	if err != nil {
		fmt.Printf("Error init: %v\n", err)
		os.Exit(1)
	}
}

func appFilePath() string {
	path, err := os.Executable()
	if err != nil {
		return os.Args[0]
	}
	return path
}
