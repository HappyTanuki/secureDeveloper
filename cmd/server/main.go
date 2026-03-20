package main

import (
	api_auth "gosecureskeleton/cmd/server/api/auth"
	"gosecureskeleton/cmd/server/api/auth/login"
	"gosecureskeleton/cmd/server/api/auth/logout"
	"gosecureskeleton/cmd/server/api/auth/register"
	auth_withdraw "gosecureskeleton/cmd/server/api/auth/withdraw"
	"gosecureskeleton/cmd/server/api/banking/deposit"
	"gosecureskeleton/cmd/server/api/banking/transfer"
	banking_withdraw "gosecureskeleton/cmd/server/api/banking/withdraw"
	"gosecureskeleton/cmd/server/api/me"
	"gosecureskeleton/cmd/server/api/posts"
	posts_id "gosecureskeleton/cmd/server/api/posts/id"
	"gosecureskeleton/cmd/server/objects/db"
	"gosecureskeleton/cmd/server/util"

	"github.com/gin-gonic/gin"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	_ "modernc.org/sqlite"

	log "github.com/sirupsen/logrus"
)

func initLogger() {
	log.SetOutput(&lumberjack.Logger{Filename: "./logs/api.log", MaxSize: 10, MaxBackups: 5, MaxAge: 30, Compress: true})
}

func JSONLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(log.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
		}).Info("Incoming Request")
		c.Next()
	}
}

func main() {
	initLogger()

	err := db.DB.OpenStore("./app.db", "./schema.sql")
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	router := gin.Default()
	util.RegisterStaticRoutes(router)
	router.Use(JSONLogger())

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", register.POST)
		auth.POST("/login", login.POST)
		auth.POST("/logout", logout.POST)
		auth.POST("/withdraw", auth_withdraw.POST)
	}

	protected := router.Group("/api")
	protected.Use(api_auth.CheckAuthority())
	{
		protected.GET("/me", me.GET)
		protected.POST("/banking/deposit", deposit.POST)
		protected.POST("/banking/withdraw", banking_withdraw.POST)
		protected.POST("/banking/transfer", transfer.POST)

		protected.GET("/posts", posts.GET)
		protected.POST("/posts", posts.POST)
		protected.GET("/posts/:id", posts_id.GET)
		protected.PUT("/posts/:id", posts_id.PUT)
		protected.DELETE("/posts/:id", posts_id.DELETE)
	}

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
