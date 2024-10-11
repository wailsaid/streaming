package controles

import (
	"encoding/base64"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/models"
	"github.com/saidwail/streaming/utils"
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
	description := c.PostForm("description")
	videoFile, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get video file"})
		return
	}

	thumbnail, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get thumbnail file"})
		return
	}

	videocodex := base64.StdEncoding.EncodeToString([]byte(videoFile.Filename))
	videoPath := filepath.Join("uploads", videocodex)
	thumbnailPath := filepath.Join("assets/thumbnails", thumbnail.Filename)

	if err := c.SaveUploadedFile(videoFile, videoPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save video file"})
		return
	}
	if err := c.SaveUploadedFile(thumbnail, thumbnailPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save thumbnail file"})
		return
	}

	// Scan video for adult content
	/* 	adultContentTimestamps, err := utils.ScanVideoForAdultContent(videoPath)
	   	if err != nil {
	   		log.Printf("Error scanning video: %v", err)
	   		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning video"})
	   		return
	   	}
	*/
	u := &models.Video{
		Title:         title,
		Description:   description,
		VideoPath:     videoPath,
		ThumbnailPath: thumbnailPath,
	}

	err = database.CreateVideo(u)
	if err != nil {
		c.Redirect(302, "/upload&s=err")
		return
	}

	c.Redirect(302, "/upload&s=ok")
}

func ListVideos(c *gin.Context) {
	list := database.GetAllVideos()

	c.JSON(http.StatusOK, list)
}

// New methods to add:

func HomePage(c *gin.Context) {
	videos := database.GetAllVideos()

	c.HTML(200, "index.html", gin.H{
		"videos": videos,
	})
}

func WatchVideo(c *gin.Context) {
	v := c.Query("v")
	video, err := database.FindVideoByID(v)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Video not found"})
		return
	}
	c.HTML(http.StatusOK, "watch.html", gin.H{"video": video})
}

func StreamVideo(c *gin.Context) {
	v := c.Query("v")
	c.File(v)
}

func RemoveAdultContent(c *gin.Context) {
	videoID := c.Param("id")
	timestamps := c.PostFormArray("timestamps")

	// Implement logic to remove adult content at specified timestamps
	err := utils.RemoveAdultContent(videoID, timestamps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove adult content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Adult content removed successfully"})
}
