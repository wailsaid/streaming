package main

import (
	"github.com/gin-gonic/gin"
	"github.com/saidwail/streaming/controles"
	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/env"
	"github.com/saidwail/streaming/utils"
)

func main() {
	env.Init()
	database.Init()

	database.Connect()

	server := gin.Default()

	server.ForwardedByClientIP = true
	server.SetTrustedProxies([]string{"127.0.0.1"})

	server.Static("./assets", "./templ/assets/**/*")

	server.LoadHTMLGlob("templ/*.html")
	server.MaxMultipartMemory = 100 << 20

	server.GET("/", controles.HomePage)
	server.GET("/watch", controles.WatchVideo)
	server.GET("/stream", controles.StreamVideo)

	server.GET("/login", controles.LoginPage)
	server.GET("/signup", controles.SignupPage)

	server.POST("/login", controles.Login)
	server.POST("/signup", controles.SignUp)

	server.GET("/upload", controles.UploadPage)
	server.POST("/upload", controles.UploadVideo)

	server.GET("/video-list", controles.ListVideos)

	server.GET("/thumbnail", controles.ServeThumbnail)

	if err := utils.InitMinioClient(); err != nil {
		//log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	server.Run()
}
