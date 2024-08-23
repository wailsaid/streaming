package controles

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/env"
	"github.com/saidwail/streaming/models"
)

func UploadPage(c *gin.Context) {
	var msg string
	switch status := c.Query("s"); status {
	case "ok":
		msg = "video uploaded successfuly"
	case "err":
		msg = "could not upload the video"
	}
	c.HTML(200, "upload.html", gin.H{
		"msg": msg,
	})
}

func UploadVideo(c *gin.Context) {
	title := c.PostForm("title")
	videoFile, err := c.FormFile("video")
	if err != nil {
		c.Redirect(http.StatusFound, "/upload?s=err")
		return
	}
	videoPath := filepath.Join("uploads", videoFile.Filename)
	u := &models.Video{
		Title: title,
		Path:  videoPath,
	}

	res := database.DB.Create(u)
	if res.Error != nil {
		log.Println(res.Error.Error())
	}

	if err := c.SaveUploadedFile(videoFile, videoPath); err != nil {
		c.Redirect(http.StatusFound, "/upload?s=err")
		return
	}
	thumbnailPath := filepath.Join("thumbnails", videoFile.Filename+".png")
	if err := env.GenerateThumbnail(videoPath, thumbnailPath); err != nil {
		c.Redirect(http.StatusFound, "/upload?s=err")
		return
	}

	c.Redirect(http.StatusFound, "/upload?s=ok")
}

func ListVideos(c *gin.Context) {
	var list []models.Video
	res := database.DB.Find(&list)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	c.JSON(http.StatusOK, list)
}
