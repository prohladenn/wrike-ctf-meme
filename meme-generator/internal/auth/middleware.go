package auth

import (
	"errors"
	"net/http"

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

	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken == "" || csrfToken != session.Get("csrf_token") {
		c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
		c.Abort()
		return
	}

	c.Set(UserIDKey, userID)

	c.Next()
}
