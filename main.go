package main

import (
	"github.com/gin-gonic/gin"
	"github.com/saidwail/streaming/controles"
	"github.com/saidwail/streaming/initEnv"
)

func main() {
	initEnv.Init()
	initEnv.Connect()

	server := gin.Default()

	server.Static("./assets", "./templ/assets")
	server.LoadHTMLGlob("templ/*.html")

	server.GET("/", func(c *gin.Context) {
		_, LogedIn := c.Get("logged_in")
		c.HTML(200, "index.html", gin.H{
			//"name":    "said",
			"LogedIn": LogedIn,
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

	server.GET("/upload", func(c *gin.Context) {
		_, LogedIn := c.Get("logged_in")
		c.HTML(200, "upload.html", gin.H{
			//"name":    "said",
			"LogedIn": LogedIn,
		})
	})
	server.POST("/upload")

	server.GET("/list_users", controles.ListUsers)

	server.Run()
}
