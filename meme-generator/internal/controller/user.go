package controller

import (
	"net/http"
	"strconv"

	"meme-generator/internal/auth"
	"meme-generator/internal/model"
	"meme-generator/internal/storage"
	"meme-generator/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserRequest struct {
	ID string `uri:"id" binding:"required"`
}

func ListUsers(s storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "100")
		pageStr := c.DefaultQuery("page", "1")

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		users, err := s.Users().All(limit, page)
		if err != nil {
			log.Error("failed to get users from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested users", zap.Int("users_count", len(users)))

		c.JSON(http.StatusOK, model.UsersToDTOs(users))
	}
}

func GetUser(s storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input UserRequest
		if err := c.ShouldBindUri(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := strconv.Atoi(input.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		user, err := s.Users().FindByID(uint(userID))
		if err != nil {
			log.Error("failed to get user from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested user info", zap.Uint("user_id", user.ID))

		c.JSON(http.StatusOK, user.ToDTO())
	}
}

func GetUserSelf(s storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet(auth.UserIDKey).(uint)

		user, err := s.Users().FindByID(userID)
		if err != nil {
			log.Error("failed to get user from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested self info", zap.Uint("user_id", user.ID))

		c.JSON(http.StatusOK, user.ToDTO())
	}
}

func ListUserTemplates(s storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input UserRequest
		if err := c.ShouldBindUri(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := strconv.Atoi(input.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		templates, err := s.Templates().FindByUserID(uint(userID))
		if err != nil {
			log.Error("failed to get templates from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested users' templates", zap.Int("templates_count", len(templates)), zap.Int("user_id", userID))

		c.JSON(http.StatusOK, model.TemplatesToDTOs(templates))
	}
}

func ListUserMemes(s storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input UserRequest
		if err := c.ShouldBindUri(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := strconv.Atoi(input.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidIntegerInput.Error()})
			return
		}

		memes, err := s.Memes().FindByUserID(uint(userID))
		if err != nil {
			log.Error("failed to get memes from the database", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User requested users' memes", zap.Int("memes_count", len(memes)), zap.Int("user_id", userID))

		c.JSON(http.StatusOK, model.MemesToDTOs(memes))
	}
}
