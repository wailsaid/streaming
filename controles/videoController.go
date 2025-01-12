package controles

import (
	"context"
	"encoding/base64"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/models"
	"github.com/saidwail/streaming/utils"
	"github.com/minio/minio-go/v7"
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

func uploadFileToMinio(file *multipart.FileHeader, bucket string) (string, error) {
	ctx := context.Background()
	minioClient := utils.GetMinioClient()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	objectName := base64.StdEncoding.EncodeToString([]byte(file.Filename))
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = minioClient.PutObject(ctx, bucket, objectName, src, file.Size,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}

	return objectName, nil
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

	// Upload files to MinIO
	videoPath, err := uploadFileToMinio(videoFile, "videos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upload video file"})
		return
	}

	thumbnailPath, err := uploadFileToMinio(thumbnail, "thumbnails")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upload thumbnail file"})
		return
	}

	u := &models.Video{
		Title:         title,
		Description:   description,
		VideoPath:     videoPath,
		ThumbnailPath: thumbnailPath,
	}

	err = database.CreateVideo(u)
	if err != nil {
		c.Redirect(302, "/upload?s=err")
		return
	}

	c.Redirect(302, "/upload?s=ok")
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
	minioClient := utils.GetMinioClient()
	
	object, err := minioClient.GetObject(context.Background(), "videos", v, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}
	defer object.Close()

	c.DataFromReader(http.StatusOK, -1, "video/mp4", object, nil)
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
