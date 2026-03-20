package login

import (
	"crypto/sha256"
	"encoding/hex"
	"gosecureskeleton/cmd/server/consts"
	"gosecureskeleton/cmd/server/objects"
	"gosecureskeleton/cmd/server/objects/db"
	"gosecureskeleton/cmd/server/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

func POST(c *gin.Context) {
	var request objects.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid login request"})
		return
	}

	user, ok, err := db.DB.FindUserByUsername(request.Username)
	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	hashed := sha256.Sum256(append([]byte(request.Password), salt...))

	if err != nil || !ok || user.Password != hex.EncodeToString(hashed[:]) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no such user or invalid credentials"})
		return
	}

	token, err := session.Sessions.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create session"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(consts.AuthorizationCookieName, token, 60*60*8, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged in"})
}
