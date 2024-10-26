package storage

import (
	"fmt"

	"meme-generator/internal/config"
	"meme-generator/internal/model"
	"meme-generator/internal/storage/meme"
	"meme-generator/internal/storage/template"
	"meme-generator/internal/storage/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

type Database struct {
	db *gorm.DB

	users     UserRepository
	templates TemplateRepository
	memes     MemeRepository
}

func NewDatabase(cfg config.DatabaseConfiguration, logger zapgorm2.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", cfg.Host, cfg.Username, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	logger.LogLevel = gormlogger.Warn
	logger.IgnoreRecordNotFoundError = true

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open postgresql connection: %w", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Template{}, &model.Meme{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func NewStorage(db *gorm.DB) *Database {
	return &Database{
		db:        db,
		users:     user.NewRepository(db),
		templates: template.NewRepository(db),
		memes:     meme.NewRepository(db),
	}
}

func (db *Database) Users() UserRepository {
	return db.users
}

func (db *Database) Templates() TemplateRepository {
	return db.templates
}

func (db *Database) Memes() MemeRepository {
	return db.memes
}
