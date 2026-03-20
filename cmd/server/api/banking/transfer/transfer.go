package transfer

import (
	"gosecureskeleton/cmd/server/objects"
	"gosecureskeleton/cmd/server/objects/db"
	"gosecureskeleton/cmd/server/session"
	"gosecureskeleton/cmd/server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func POST(c *gin.Context) {
	var request objects.TransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid transfer request"})
		return
	}
	// already checked from middleware
	token := util.TokenFromRequest(c)
	sessionData, _ := session.Sessions.Lookup(token)

	userB, ok, err := db.DB.FindUserByUsername(request.ToUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "target user doesn't exist"})
		return
	}

	ok, err = db.DB.TransferBalenceAToB(sessionData.User.ID, userB.ID, request.Amount)
	if !ok && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "insufficient balence"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "transferred",
		"target":  request.ToUsername,
		"amount":  request.Amount,
	})
}
