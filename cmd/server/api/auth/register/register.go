package register

import (
	"crypto/sha256"
	"encoding/hex"
	"gosecureskeleton/cmd/server/objects"
	"gosecureskeleton/cmd/server/objects/db"
	"gosecureskeleton/cmd/server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func POST(c *gin.Context) {
	var request objects.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid register request"})
		return
	}

	if _, ok, err := db.DB.FindUserByUsername(request.Username); ok || err != nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusConflict, gin.H{"message": request.Username + "is already occupied"})
		return
	}

	salt, err := util.GenerateRandomBytes(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
	hashed := sha256.Sum256(append([]byte(request.Password), salt...))

	err = db.DB.InsertUser(objects.User{
		Username: request.Username,
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: hex.EncodeToString(hashed[:]),
		Salt:     hex.EncodeToString(salt),
		IsAdmin:  false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Accepted",
	})
}
