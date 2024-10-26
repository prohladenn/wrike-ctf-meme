package controller

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"meme-generator/internal/auth"
	"meme-generator/internal/model"
	"meme-generator/internal/storage"
	"meme-generator/log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"golang.org/x/time/rate"
)

const (
	minCredsLength    = 9
	maxUsernameLength = 16
	maxPasswordLength = 32
)

var (
	errUsernameExists          = errors.New("username already exists")
	errEmptyUsernameOrPassword = errors.New("username or password is empty")
	errCredsLength             = errors.New("username or password length is incorrect")
	errHashPassword            = errors.New("failed to hash password")
	errInvalidCredentials      = errors.New("invalid credentials")
	errInvalidUsernameSymbols  = errors.New("invalid username symbols: only latin letters, underscores, dashes and numbers are allowed")
	errTooManyRequests         = errors.New("too many requests, please try again later")
)

const (
	msgRegistered = "registered successfully"
	msgLoggedIn   = "logged in successfully"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var limiter = rate.NewLimiter(1, 5)

func Register(s storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": errTooManyRequests.Error()})
			return
		}

		var input AuthRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if ok, err := s.Users().CheckExistsByUsername(input.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": errUsernameExists.Error()})
			return
		}

		if err := checkAuthRequest(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := auth.HashPassword(input.Password)
		if err != nil {
			log.Error("failed to hash password", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errHashPassword.Error()})
			return
		}

		user := model.User{
			Username: input.Username,
			Password: string(hashedPassword),
		}

		if _, err := s.Users().Create(&user); err != nil {
			log.Error("failed to create user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		log.Info("User registered", zap.Uint("user_id", user.ID), zap.String("username", user.Username))

		c.JSON(http.StatusCreated, gin.H{"message": msgRegistered})
	}
}

func Login(s storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": errTooManyRequests.Error()})
			return
		}

		var input AuthRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := s.Users().FindByUsername(input.Username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidCredentials.Error()})
				return
			}

			log.Error("failed to find user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDatabase.Error()})
			return
		}

		if !auth.CheckPassword(user.Password, input.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidCredentials.Error()})
			return
		}

		session := sessions.Default(c)
		session.Set(auth.UserIDKey, user.ID)

		if err := session.Save(); err != nil {
			log.Error("failed to save session", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errSessionSave.Error()})
			return
		}

		auth.RegenerateSession(c)
		auth.SetSessionExpiration(c, time.Hour)

		log.Info("User logged in", zap.Uint("user_id", user.ID), zap.String("username", user.Username))

		c.JSON(http.StatusOK, gin.H{"message": msgLoggedIn})
	}
}

func Logout(s storage.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete(auth.UserIDKey)
		session.Save()

		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
	}
}

func checkAuthRequest(input *AuthRequest) error {
	if input.Username == "" || input.Password == "" {
		return errEmptyUsernameOrPassword
	}

	if len(input.Username) <= minCredsLength || len(input.Username) >= maxUsernameLength || len(input.Password) <= minCredsLength || len(input.Password) >= maxPasswordLength {
		return errCredsLength
	}

	regexp := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if !regexp.MatchString(input.Username) {
		return errInvalidUsernameSymbols
	}

	return nil
}
