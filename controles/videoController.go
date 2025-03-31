package controles

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"text/template"

	"github.com/minio/minio-go/v7"
	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/models"
	"github.com/saidwail/streaming/utils"
)

// Map is a shorthand for map[string]interface{}
type Map map[string]interface{}

func UploadPage(c *CustomContext) {
	if c.GetHeader("HX-Request") == "true" {
		// For HTMX requests, only render the content portion
		c.PartialHTML(http.StatusOK, "upload", "content", Map{
			"title": "upload",
		})

		// Update the browser title via HTMX
		c.Writer.Header().Add("HX-Push-Url", "/")
		c.Writer.Header().Add("HX-Trigger", `{"updateTitle": "Uploade - YouClone"}`)

		return
	}

	var msg string
	switch status := c.Query("s"); status {
	case "ok":
		msg = "video uploaded successfuly"
	case "err":
		msg = "could not upload the video"
	}
	c.HTML(200, "upload", Map{
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

func UploadVideo(c *CustomContext) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	videoFile, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{"error": "Could not get video file"})
		return
	}

	thumbnail, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{"error": "Could not get thumbnail file"})
		return
	}

	// Upload files to MinIO
	videoPath, err := uploadFileToMinio(videoFile, "videos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "Could not upload video file"})
		return
	}

	thumbnailPath, err := uploadFileToMinio(thumbnail, "thumbnails")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "Could not upload thumbnail file"})
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

func ListVideos(c *CustomContext) {
	list := database.GetAllVideos()

	c.JSON(http.StatusOK, list)
}

// New methods to add:

func HomePage(c *CustomContext) {
	videos := database.GetAllVideos()

	// Check if this is an HTMX request
	if c.GetHeader("HX-Request") == "true" {
		// For HTMX requests, only render the content portion
		c.PartialHTML(http.StatusOK, "index", "content", Map{
			"videos": videos,
			"title":  "home",
		})

		// Update the browser title via HTMX
		c.Header("HX-Push-Url", "/")
		c.Header("HX-Trigger", `{"updateTitle": "Home - YouClone"}`)
		return
	}

	c.HTML(200, "index", Map{
		"videos": videos,
		"title":  "home",
	})
}

func WatchVideo(c *CustomContext) {
	v := c.Query("v")
	video, err := database.FindVideoByID(v)
	if err != nil {
		c.HTML(http.StatusNotFound, "error", Map{"error": "Video not found"})
		return
	}

	// Get recommended videos (excluding current video)
	recommendations := database.GetRecommendedVideos(v, 10) // Get 10 recommendations

	// Get video comments
	comments := database.GetCommentsByVideoID(video.ID, 5) // Get 5 latest comments

	// Check if this is an HTMX request
	if c.GetHeader("HX-Request") == "true" {
		// For HTMX requests, only render the content portion
		c.PartialHTML(http.StatusOK, "watch", "content", Map{
			"video":           video,
			"recommendations": recommendations,
			"comments":        comments,
			"hideSide":        true,
		})

		// Update the browser title via HTMX
		c.Header("HX-Push-Url", "/watch?v="+v)
		c.Header("HX-Trigger", fmt.Sprintf(`{"updateTitle": "%s - YouClone"}`, video.Title))
		return
	}

	// For regular requests, render the full page
	c.HTML(http.StatusOK, "watch", Map{
		"video":           video,
		"recommendations": recommendations,
		"comments":        comments,
		"hideSide":        true,
	})
}

func StreamVideo(c *CustomContext) {
	v := c.Query("v")
	minioClient := utils.GetMinioClient()

	// Get object info to get the size
	objInfo, err := minioClient.StatObject(context.Background(), "videos", v, minio.StatObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, Map{"error": "Video not found"})
		return
	}

	// Parse Range header
	rangeHeader := c.GetHeader("Range")
	start, end := int64(0), objInfo.Size-1

	if rangeHeader != "" {
		// Parse the Range header value
		_, err := fmt.Sscanf(rangeHeader, "bytes=%d-", &start)
		if err != nil {
			start = 0
		}

		if start < 0 {
			start = 0
		}

		if start >= objInfo.Size {
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}
	}

	// Get object with range
	opts := minio.GetObjectOptions{}
	if err := opts.SetRange(start, end); err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "Invalid range"})
		return
	}

	object, err := minioClient.GetObject(context.Background(), "videos", v, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "Could not retrieve video"})
		return
	}
	defer object.Close()

	// Set the appropriate headers
	c.Header("Content-Type", "video/mp4")
	c.Header("Accept-Ranges", "bytes")
	contentLength := end - start + 1
	contentRange := fmt.Sprintf("bytes %d-%d/%d", start, end, objInfo.Size)

	if rangeHeader != "" {
		c.Header("Content-Range", contentRange)
		c.Status(http.StatusPartialContent)
	} else {
		c.Status(http.StatusOK)
	}

	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))

	// Stream the video
	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := object.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if n > 0 {
			c.Writer.Write(buffer[:n])
			c.Writer.(http.Flusher).Flush()
		}
	}
}

func RemoveAdultContent(c *CustomContext) {
	videoID := c.Param("id")
	timestamps := c.PostFormArray("timestamps")

	// Implement logic to remove adult content at specified timestamps
	err := utils.RemoveAdultContent(videoID, timestamps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "Failed to remove adult content"})
		return
	}

	c.JSON(http.StatusOK, Map{"message": "Adult content removed successfully"})
}

// LoadMoreVideos handles HTMX requests to load more videos
func LoadMoreVideos(c *CustomContext) {
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 8 // Default limit
	}

	// Get videos with pagination
	videos := database.GetPaginatedVideos(offset, limit)

	// Instead of rendering a full HTML page, render only the video links component
	tmplFiles := []string{
		"templates/components/vlink.html",
	}

	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(c.Writer, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(http.StatusOK)

	for _, video := range videos {
		err = t.ExecuteTemplate(c.Writer, "vlink", video)
		if err != nil {
			http.Error(c.Writer, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// SearchVideos handles HTMX requests to search videos
func SearchVideos(c *CustomContext) {
	query := c.Query("q")
	if query == "" {
		// If no query provided, return all videos
		videos := database.GetAllVideos()
		renderVideoResults(c, videos)
		return
	}

	// Search videos by title or description
	videos := database.SearchVideos(query)
	renderVideoResults(c, videos)
}

// Helper function to render video search results
func renderVideoResults(c *CustomContext, videos []models.Video) {
	tmplFiles := []string{
		"templates/components/vlink.html",
	}

	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(c.Writer, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(http.StatusOK)

	if len(videos) == 0 {
		c.Writer.Write([]byte("<div class='col-span-full text-center py-8'>No videos found matching your search.</div>"))
		return
	}

	for _, video := range videos {
		err = t.ExecuteTemplate(c.Writer, "vlink", video)
		if err != nil {
			http.Error(c.Writer, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// VideoPreview provides a GIF or video snippet preview when hovering
func VideoPreview(c *CustomContext) {
	v := c.Query("v")
	if v == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	// In a real implementation, this would generate or serve an animated preview
	// For now, we'll just provide a simple HTML overlay with a play button

	html := `
	<div class="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center">
		<div class="animate-pulse flex flex-col items-center">
			<svg xmlns="http://www.w3.org/2000/svg" width="36" height="36" viewBox="0 0 24 24" fill="none" 
				stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-white mb-2">
				<polygon points="5 3 19 12 5 21 5 3"/>
			</svg>
			<span class="text-white text-sm font-medium">Watch Now</span>
		</div>
	</div>
	`

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write([]byte(html))
}

func ServeThumbnail(c *CustomContext) {
	v := c.Query("v")
	minioClient := utils.GetMinioClient()

	object, err := minioClient.GetObject(context.Background(), "thumbnails", v, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, Map{"error": "Thumbnail not found"})
		return
	}
	defer object.Close()

	// Get content type
	objInfo, err := minioClient.StatObject(context.Background(), "thumbnails", v, minio.StatObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "Could not get thumbnail info"})
		return
	}

	contentType := objInfo.ContentType
	if contentType == "" {
		contentType = "image/jpeg" // default content type for images
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", objInfo.Size))

	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := object.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if n > 0 {
			c.Writer.Write(buffer[:n])
		}
	}
}
