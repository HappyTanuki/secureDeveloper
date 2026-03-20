package withdraw

import (
	"gosecureskeleton/cmd/server/objects"
	"gosecureskeleton/cmd/server/objects/db"
	"gosecureskeleton/cmd/server/session"
	"gosecureskeleton/cmd/server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func POST(c *gin.Context) {
	var request objects.BalanceWithdrawRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid withdraw request"})
		return
	}

	// already checked from middleware
	token := util.TokenFromRequest(c)
	sessionData, _ := session.Sessions.Lookup(token)

	ok, err := db.DB.AddUserBalenceByID(sessionData.User.ID, request.Amount*(-1))
	if !ok && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "insufficient balence"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "withdrawed",
		"amount":  request.Amount,
	})
}
