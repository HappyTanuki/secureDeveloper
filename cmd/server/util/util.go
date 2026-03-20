package util

import (
	"crypto/rand"
	"encoding/hex"
	"gosecureskeleton/cmd/server/consts"
	"gosecureskeleton/cmd/server/objects"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterStaticRoutes(router *gin.Engine) {
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
}

func MakeUserResponse(user objects.User) objects.UserResponse {
	return objects.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Balance:  user.Balance,
		IsAdmin:  user.IsAdmin,
	}
}

func ClearAuthorizationCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(consts.AuthorizationCookieName, "", -1, "/", "", false, true)
}

func TokenFromRequest(c *gin.Context) string {
	headerValue := strings.TrimSpace(c.GetHeader("Authorization"))
	if headerValue != "" {
		return headerValue
	}

	cookieValue, err := c.Cookie(consts.AuthorizationCookieName)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(cookieValue)
}

func GenerateRandomBytes(size int) ([]byte, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return nil, err
	}
	return buffer, nil
}

func GenerateToken() (string, error) {
	var buffer []byte
	var err error
	if buffer, err = GenerateRandomBytes(24); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}
