package env

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/saidwail/streaming/database"
)

func Init() {

	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}

	if mode := os.Getenv("GIN_MODE"); mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)

	}
	/* 	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
	   		log.Fatal(err.Error())
	   	}
	   	if err := os.MkdirAll("thumbnails", os.ModePerm); err != nil {
	   		log.Fatal(err.Error())
	   	} */

	database.Init()
}
