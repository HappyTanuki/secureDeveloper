package id

import (
	"gosecureskeleton/cmd/server/objects"
	"gosecureskeleton/cmd/server/session"
	"gosecureskeleton/cmd/server/util"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GET(c *gin.Context) {
	token := util.TokenFromRequest(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization token"})
		return
	}
	if _, ok := session.Sessions.Lookup(token); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization token"})
		return
	}

	c.JSON(http.StatusOK, objects.PostResponse{
		Post: objects.PostView{
			ID:          1,
			Title:       "Dummy Post",
			Content:     "This is a fixed dummy response. Replace this later with real board logic.",
			OwnerID:     1,
			Author:      "Alice Admin",
			AuthorEmail: "alice.admin@example.com",
			CreatedAt:   "2026-03-19T09:00:00Z",
			UpdatedAt:   "2026-03-19T09:00:00Z",
		},
	})
}

func PUT(c *gin.Context) {
	var request objects.UpdatePostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid update request"})
		return
	}

	token := util.TokenFromRequest(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization token"})
		return
	}
	user, ok := session.Sessions.Lookup(token)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization token"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	c.JSON(http.StatusOK, gin.H{
		"message": "dummy update post handler",
		"todo":    "replace with ownership check and update query",
		"post": objects.PostView{
			ID:          1,
			Title:       strings.TrimSpace(request.Title),
			Content:     strings.TrimSpace(request.Content),
			OwnerID:     user.User.ID,
			Author:      user.User.Name,
			AuthorEmail: user.User.Email,
			CreatedAt:   "2026-03-19T09:00:00Z",
			UpdatedAt:   now,
		},
	})
}

func DELETE(c *gin.Context) {
	token := util.TokenFromRequest(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization token"})
		return
	}
	if _, ok := session.Sessions.Lookup(token); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "dummy delete post handler",
		"todo":    "replace with ownership check and delete query",
	})
}
