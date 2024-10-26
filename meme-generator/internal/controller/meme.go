package controller

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"meme-generator/internal/auth"
	"meme-generator/internal/meme"
	"meme-generator/internal/model"
	"meme-generator/internal/storage"
	"meme-generator/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	previewMemeName = "meme-preview"
)

var (
	errMemeNotFound = errors.New("meme not found")
	errCreateMeme   = errors.New("failed to create meme")
)

type CreateMemeRequest struct {
	TemplateID string `json:"template_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Caption    string `json:"caption" binding:"required"`
	Color      string `json:"color"`
}

type PreviewMemeRequest struct {
	TemplateID string `json:"template_id" binding:"required"`
	Caption    string `json:"caption" binding:"required"`
	Color      string `json:"color"`
}

type MemeRequest struct {
	ID string `uri:"id" binding:"required"`
}

func CreateMeme(s storage.Store, memeManager *meme.MemeManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet(auth.UserIDKey).(uint)

		var input CreateMemeRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		templateID, err := strconv.Atoi(input.TemplateID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		template, err := s.Templates().FindByID(uint(templateID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": errMemeNotFound.Error()})
				return
			}

			log.Error("failed to get template from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		templateOwner, err := s.Users().FindByID(template.UserID)
		if err != nil {
			log.Error("failed to get user from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		meme, err := memeManager.NewMeme(template, templateOwner.Username, userID, input.Name, input.Caption, input.Color)
		if err != nil {
			log.Error("failed to create meme", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errCreateMeme.Error()})
			return
		}

		createdMeme, err := s.Memes().Create(meme)
		if err != nil {
			log.Error("failed to create meme in the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User created meme", zap.Uint("meme_id", createdMeme.ID), zap.String("meme_name", createdMeme.Name), zap.String("meme_caption", input.Caption), zap.Uint("user_id", userID), zap.Int("template_id", templateID))

		c.JSON(http.StatusCreated, createdMeme.ToDTO())
	}
}

func GetMeme(s storage.Store, memeManager *meme.MemeManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input MemeRequest
		if err := c.ShouldBindUri(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := strconv.Atoi(input.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		meme, err := s.Memes().FindByID(uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": errMemeNotFound.Error()})
				return
			}

			log.Error("failed to get meme from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested meme info", zap.Uint("meme_id", meme.ID), zap.String("meme_name", meme.Name), zap.Uint("user_id", meme.UserID))

		c.JSON(http.StatusOK, meme.ToDTO())
	}
}

func GetMemeImage(s storage.Store, memeManager *meme.MemeManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input MemeRequest
		if err := c.ShouldBindUri(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := strconv.Atoi(input.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		meme, err := s.Memes().FindByID(uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": errMemeNotFound.Error()})
				return
			}

			log.Error("failed to get meme from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested meme image", zap.Uint("meme_id", meme.ID), zap.String("meme_name", meme.Name), zap.Uint("user_id", meme.UserID))

		filePath := filepath.Join(meme.DirPath, meme.FileName)
		c.FileAttachment(filePath, meme.FileName)
	}
}

func GetLastMemes(s storage.Store, memeManager *meme.MemeManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		memes, err := s.Memes().LastRecords()
		if err != nil {
			log.Error("failed to get memes from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested last memes", zap.Int("memes_count", len(memes)))

		c.JSON(http.StatusOK, model.MemesToDTOs(memes))
	}
}

func PreviewMeme(s storage.Store, memeManager *meme.MemeManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet(auth.UserIDKey).(uint)

		var input PreviewMemeRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		templateID, err := strconv.Atoi(input.TemplateID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		template, err := s.Templates().FindByID(uint(templateID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": errMemeNotFound.Error()})
				return
			}

			log.Error("failed to get template from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		templateOwner, err := s.Users().FindByID(template.UserID)
		if err != nil {
			log.Error("failed to get user from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		meme, err := memeManager.NewMeme(template, templateOwner.Username, userID, previewMemeName, input.Caption, input.Color)
		if err != nil {
			log.Error("failed to create meme", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errCreateMeme.Error()})
			return
		}

		log.Info("User requested meme preview", zap.Uint("meme_id", meme.ID), zap.String("meme_name", meme.Name), zap.Uint("user_id", meme.UserID))

		filePath := filepath.Join(meme.DirPath, meme.FileName)
		c.FileAttachment(filePath, meme.FileName)
	}
}
