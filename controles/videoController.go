package controles

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saidwail/streaming/initEnv"
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
	videoFile, err := c.FormFile("video")
	if err != nil {
		c.Redirect(http.StatusFound, "/upload?s=err")
		return
	}
	path := "./Videos/" + videoFile.Filename
	u := &models.Video{
		Title: c.Query("title"),
		Path:  path,
	}

	initEnv.DB.Create(u)

	if err := c.SaveUploadedFile(videoFile, path); err != nil {
		c.Redirect(http.StatusFound, "/upload?s=err")
		return
	}

	c.Redirect(http.StatusFound, "/upload?s=ok")
}

func ListVideos(c *gin.Context) {
	var list models.Video
	res := initEnv.DB.Find(&list)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	c.JSON(http.StatusOK, list)
}
