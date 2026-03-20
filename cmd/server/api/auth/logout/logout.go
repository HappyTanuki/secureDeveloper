package logout

import (
	"gosecureskeleton/cmd/server/session"
	"gosecureskeleton/cmd/server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func POST(c *gin.Context) {
	token := util.TokenFromRequest(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization token"})
		return
	}
	if _, ok := session.Sessions.Lookup(token); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization token"})
		return
	}

	session.Sessions.Delete(token)
	util.ClearAuthorizationCookie(c)
	c.JSON(http.StatusOK, gin.H{
		"message": "logout",
	})
}
