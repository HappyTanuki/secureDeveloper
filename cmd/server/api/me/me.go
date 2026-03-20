package me

import (
	"gosecureskeleton/cmd/server/objects/db"
	"gosecureskeleton/cmd/server/session"
	"gosecureskeleton/cmd/server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GET(c *gin.Context) {
	// already checked from middleware
	token := util.TokenFromRequest(c)
	sessionData, _ := session.Sessions.Lookup(token)

	user, ok, err := db.DB.FindUserByUsername(sessionData.User.Username)
	if !ok && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": util.MakeUserResponse(user)})
}
