package controles

import (
	"encoding/base64"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/saidwail/streaming/database"
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

	thumbnail, err := c.FormFile("thumbnail")
	if err != nil {
		c.Redirect(http.StatusFound, "/upload?s=err")
		return
	}
	videocodex := base64.StdEncoding.EncodeToString([]byte(videoFile.Filename))

	videoPath := filepath.Join("uploads", videocodex)
	thumbnailPath := filepath.Join("assets", thumbnail.Filename)

	u := &models.Video{
		Title:         title,
		VideoPath:     videoPath,
		ThumbnailPath: thumbnailPath,
	}

	res := database.DB.Create(u)
	if res.Error != nil {
		log.Fatal(res.Error.Error())
	}

	if err := c.SaveUploadedFile(videoFile, videoPath); err != nil {
		c.Redirect(http.StatusFound, "/upload?s=err")
		return
	}
	if err := c.SaveUploadedFile(thumbnail, thumbnailPath); err != nil {
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
