package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var errUnauthorized = errors.New("unauthorized")

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get(UserIDKey)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errUnauthorized.Error()})
		c.Abort()
		return
	}

	c.Set(UserIDKey, userID)

	c.Next()
}

func RegenerateSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		Secure:   true,
		HttpOnly: true,
	})
	session.Save()
}

func SetSessionExpiration(c *gin.Context, duration time.Duration) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   int(duration.Seconds()),
		Secure:   true,
		HttpOnly: true,
	})
	session.Save()
}
