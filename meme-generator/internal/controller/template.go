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

var (
	errTemplateNotFound  = errors.New("template not found")
	errEmptyFile         = errors.New("empty file")
	errCreateTemplate    = errors.New("failed to create template")
	errEmptyTemplateName = errors.New("empty template name")
	errCheckFileType     = errors.New("failed to check file type")
	errInvalidFileType   = errors.New("invalid file type: only jpeg, jpg and png are allowed")
	errTemplateInfo      = errors.New("failed to get template private info")
	errFileTooLarge      = errors.New("file size exceeds the limit of 5MB")
)

const maxFileSize = 5 * 1024 * 1024 // 5MB

type TemplateRequest struct {
	ID string `uri:"id" binding:"required"`
}

func CreateTemplate(s storage.Store, memeManager *meme.MemeManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet(auth.UserIDKey).(uint)

		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if file.Size == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": errEmptyFile.Error()})
			return
			}
		if file.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": errFileTooLarge.Error()})
			return
		}

		ok, err := meme.CheckImageType(file)
		if err != nil {
			log.Error("failed to check file type", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errCheckFileType})
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidFileType.Error()})
			return
		}

		name := c.PostForm("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errEmptyTemplateName.Error()})
			return
		}

		comment := c.PostForm("comment")

		template, err := memeManager.NewTemplate(name, comment, file, userID)
		if err != nil {
			log.Error("failed to prepare template directory", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errCreateTemplate.Error()})
			return
		}

		createdTemplate, err := s.Templates().Create(template)
		if err != nil {
			log.Error("failed to create template in the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		privateInfo, err := memeManager.GetTemplatePrivateInfo(createdTemplate)
		if err != nil {
			log.Error("failed to get template private info", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errTemplateInfo.Error()})
			return
		}

		log.Info("User created template", zap.Uint("template_id", createdTemplate.ID), zap.String("template_name", createdTemplate.Name), zap.String("template_comment", comment), zap.Uint("user_id", userID))

		c.JSON(http.StatusOK, createdTemplate.ToPrivateDTO(privateInfo))
	}
}

func GetTemplate(s storage.Store, memeManager *meme.MemeManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input TemplateRequest
		if err := c.ShouldBindUri(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := strconv.Atoi(input.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		template, err := s.Templates().FindByID(uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": errTemplateNotFound.Error()})
				return
			}

			log.Error("failed to get template from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		userID := c.MustGet(auth.UserIDKey).(uint)

		if template.UserID == userID {
			privateInfo, err := memeManager.GetTemplatePrivateInfo(template)
			if err != nil {
				log.Error("failed to get template private info", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": errTemplateInfo.Error()})
				return
			}

			c.JSON(http.StatusOK, template.ToPrivateDTO(privateInfo))
			return
		}

		log.Info("User requested template", zap.Uint("template_id", template.ID), zap.String("template_name", template.Name), zap.Uint("user_id", userID))

		c.JSON(http.StatusOK, template.ToDTO())
	}
}

func GetTemplateImage(s storage.Store, memeManager *meme.MemeManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input TemplateRequest
		if err := c.ShouldBindUri(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := strconv.Atoi(input.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		template, err := s.Templates().FindByID(uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": errTemplateNotFound.Error()})
				return
			}

			log.Error("failed to get template from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested template image", zap.Uint("template_id", template.ID), zap.String("template_name", template.Name), zap.Uint("user_id", template.UserID))

		filePath := filepath.Join(template.DirPath, template.FileName)
		c.FileAttachment(filePath, template.FileName)
	}
}

func GetLastTemplates(s storage.Store, memeManager *meme.MemeManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		templates, err := s.Templates().LastRecords()
		if err != nil {
			log.Error("failed to get last templates from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested last templates", zap.Int("templates_count", len(templates)))

		c.JSON(http.StatusOK, model.TemplatesToDTOs(templates))
	}
}
