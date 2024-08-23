package main

import (
	"github.com/gin-gonic/gin"
	"github.com/saidwail/streaming/controles"
	"github.com/saidwail/streaming/initEnv"
	"github.com/saidwail/streaming/models"
)

func main() {
	initEnv.Init()
	initEnv.Connect()

	server := gin.Default()

	server.Static("./assets", "./templ/assets")
	server.LoadHTMLGlob("templ/*.html")
	server.MaxMultipartMemory = 100 << 20

	server.GET("/", func(c *gin.Context) {
		var videos []models.Video
		initEnv.DB.Find(&videos)

		c.HTML(200, "index.html", gin.H{
			"videos": videos,
		})
	})

	server.GET("signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", nil)
	})
	server.POST("/signup", controles.SignUp)

	server.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	server.POST("/login", controles.Login)

	server.GET("/upload", controles.UploadPage)

	server.POST("/upload", controles.UploadVideo)

	server.GET("/video-list", controles.ListVideos)
	//server.GET("/list_users", controles.ListUsers)

	server.Run()
}
