package main

import (
	"github.com/gin-gonic/gin"
	"github.com/saidwail/learningGo/controles"
	"github.com/saidwail/learningGo/initEnv"
	"github.com/saidwail/learningGo/midelware"
)

func main() {
	initEnv.Init()
	initEnv.Connect()

	server := gin.Default()

	server.Static("/assets", "./templ/assets")
	server.LoadHTMLGlob("templ/*.html")

	server.GET("/", midelware.JwtFilter, func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"name": "<h1>said</h1>",
		})
	})

	server.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	server.GET("signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", nil)
	})

	server.POST("/signup", controles.AddUser)

	server.POST("/login", controles.Login)

	server.GET("/list_users", controles.ListUsers)

	server.Run()
}
