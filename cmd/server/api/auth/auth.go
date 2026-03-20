package auth

import (
	"gosecureskeleton/cmd/server/session"
	"gosecureskeleton/cmd/server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuthority() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := util.TokenFromRequest(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization token"})
			c.Abort()
		}
		_, ok := session.Sessions.Lookup(token)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization token"})
			c.Abort()
		}
		c.Next()
	}
}
