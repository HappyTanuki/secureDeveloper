package withdraw

import (
	"crypto/sha256"
	"encoding/hex"
	"gosecureskeleton/cmd/server/objects"
	"gosecureskeleton/cmd/server/objects/db"
	"gosecureskeleton/cmd/server/session"
	"gosecureskeleton/cmd/server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func POST(c *gin.Context) {
	var request objects.WithdrawAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid withdraw request"})
		return
	}

	token := util.TokenFromRequest(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization token"})
		return
	}
	sessionData, ok := session.Sessions.Lookup(token)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization token"})
		return
	}

	salt, err := hex.DecodeString(sessionData.User.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	hashed := sha256.Sum256(append([]byte(request.Password), salt...))

	if sessionData.User.Password != hex.EncodeToString(hashed[:]) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect password"})
		return
	}

	db.DB.DeleteUserByID(sessionData.User.ID)

	c.JSON(http.StatusAccepted, gin.H{"message": "User " + sessionData.User.Username + " deleted."})
}
