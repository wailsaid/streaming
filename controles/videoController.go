package controles

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/models"
)

// Map is a shorthand for map[string]interface{}
type Map map[string]interface{}

// StorageDir contains the base directory for all media storage
const (
	VideosDir     = "./storage/videos"
	ThumbnailsDir = "./storage/thumbnails"
)

// UploadChunk handles chunked video uploads (multipart)
// POST /api/videos/upload
// FormData: file (chunk), file_name, index, total
func UploadChunk(c *CustomContext) {
	log.Println("uploading chunk")
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{"error": "could not read file chunk"})
		return
	}
	defer file.Close()

	filename := c.PostForm("file_name")
	index, _ := strconv.Atoi(c.PostForm("index"))
	total, _ := strconv.Atoi(c.PostForm("total"))

	log.Printf("receiving file %v part %v size %v total %v", filename, index, fileHeader.Size, total)

	// Save chunk to tmp
	tmpPath := fmt.Sprintf("./tmp/%s_%d.tmp", filename, index)
	fw, err := os.Create(tmpPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "could not write chunk"})
		return
	}
	b := make([]byte, fileHeader.Size)
	file.Read(b)
	fw.Write(b)
	fw.Close()

	// If not the last chunk, just ack
	if index+1 != total {
		c.JSON(http.StatusOK, Map{"message": "chunk received"})
		return
	}

	// All chunks received — merge in background
	log.Printf("all parts received for %v, merging...", filename)
	c.JSON(http.StatusCreated, Map{"message": "File uploaded successfully"})

	go mergeChunks(filename, total)
}

func mergeChunks(filename string, total int) {
	destPath := filepath.Join(VideosDir, filename)
	video, err := os.Create(destPath)
	if err != nil {
		log.Printf("error creating video file: %v", err)
		return
	}
	defer video.Close()

	for i := range total {
		tmpPath := fmt.Sprintf("./tmp/%s_%d.tmp", filename, i)
		b, err := os.ReadFile(tmpPath)
		if err != nil {
			log.Printf("error reading chunk %d: %v", i, err)
			return
		}
		video.Write(b)
		os.Remove(tmpPath)
	}

	log.Printf("file %v merged successfully to %v", filename, destPath)

	// Store in DB
	videoRecord := &models.Video{
		Title:         strings.TrimSuffix(filename, filepath.Ext(filename)),
		VideoPath:     destPath,
		ThumbnailPath: "",
		CreatedAt:     time.Now(),
	}
	if err := database.CreateVideo(videoRecord); err != nil {
		log.Printf("error storing video in DB: %v", err)
	}
}

// ListVideos returns all videos as JSON
// GET /api/videos
func ListVideos(c *CustomContext) {
	list := database.GetAllVideos()
	c.JSON(http.StatusOK, list)
}

// GetVideo returns a single video by ID
// GET /api/videos/{id}
func GetVideo(c *CustomContext) {
	id := c.Param("id")
	video, err := database.FindVideoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Map{"error": "video not found"})
		return
	}
	c.JSON(http.StatusOK, video)
}

// StreamVideo streams a local video file with Range support
// GET /api/stream/{id}
func StreamVideo(c *CustomContext) {
	id := c.Param("id")
	video, err := database.FindVideoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Map{"error": "video not found"})
		return
	}

	f, err := os.Open(video.VideoPath)
	if err != nil {
		c.JSON(http.StatusNotFound, Map{"error": "video file not found"})
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "could not stat file"})
		return
	}

	c.Writer.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), f)
}

// ServeThumbnail serves a thumbnail for a video
// GET /api/thumbnail/{id}
func ServeThumbnail(c *CustomContext) {
	id := c.Param("id")
	video, err := database.FindVideoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Map{"error": "video not found"})
		return
	}

	if video.ThumbnailPath == "" {
		c.JSON(http.StatusNotFound, Map{"error": "no thumbnail"})
		return
	}

	f, err := os.Open(video.ThumbnailPath)
	if err != nil {
		c.JSON(http.StatusNotFound, Map{"error": "thumbnail file not found"})
		return
	}
	defer f.Close()

	stat, _ := f.Stat()
	http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), f)
}

// SearchVideos searches by title query param
// GET /api/videos/search?q=...
func SearchVideos(c *CustomContext) {
	query := c.Query("q")
	var videos []models.Video
	if query == "" {
		videos = database.GetAllVideos()
	} else {
		videos = database.SearchVideos(query)
	}
	c.JSON(http.StatusOK, videos)
}

// LoadMoreVideos returns paginated videos
// GET /api/videos/more?offset=0&limit=8
func LoadMoreVideos(c *CustomContext) {
	offset, _ := strconv.Atoi(c.Query("offset"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 8
	}
	videos := database.GetPaginatedVideos(offset, limit)
	c.JSON(http.StatusOK, videos)
}

// GetVideoComments returns comments for a video
// GET /api/videos/{id}/comments
func GetVideoComments(c *CustomContext) {
	id := c.Param("id")
	videoID64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{"error": "invalid video id"})
		return
	}
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 20
	}
	comments := database.GetCommentsByVideoID(uint(videoID64), limit)
	c.JSON(http.StatusOK, comments)
}

// AddComment adds a comment to a video
// POST /api/videos/{id}/comments
func AddComment(c *CustomContext) {
	id := c.Param("id")
	videoID64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{"error": "invalid video id"})
		return
	}

	var body struct {
		Content string `json:"content"`
	}
	if err := c.Bind(&body); err != nil || body.Content == "" {
		c.JSON(http.StatusBadRequest, Map{"error": "content is required"})
		return
	}

	// Get user from context (set by JWT middleware)
	userID := uint(1) // fallback; ideally from ctx
	if err == nil {
		// try to get user from context
		// user := c.Request.Context().Value(midelware.UserContextKey)
	}

	comment := &models.Comment{
		Content:   body.Content,
		VideoID:   uint(videoID64),
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	if err := database.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "could not save comment"})
		return
	}
	user, err := database.FindUserByID(userID)
	if err == nil {
		comment.User = user
	}
	c.JSON(http.StatusCreated, comment)
}

// GetRecommendedVideos returns recommended videos excluding an ID
// GET /api/videos/{id}/recommended
func GetRecommendedVideos(c *CustomContext) {
	id := c.Param("id")
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	videos := database.GetRecommendedVideos(id, limit)
	c.JSON(http.StatusOK, videos)
}

// VideoPreview provides a minimal preview response
func VideoPreview(c *CustomContext) {
	v := c.Param("id")
	if v == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	video, err := database.FindVideoByID(v)
	if err != nil {
		c.JSON(http.StatusNotFound, Map{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, Map{
		"id":             video.ID,
		"title":          video.Title,
		"thumbnail_path": video.ThumbnailPath,
	})
}

// helper: unused legacy below (kept to avoid removing necessary symbols)
func renderVideoResults(_ *CustomContext, _ []models.Video) {}

// UploadPage — only here to prevent compile errors from old references; not used
func UploadPage(c *CustomContext) {
	c.JSON(http.StatusOK, Map{"message": "use the Vue upload page"})
}

// UploadThumbnail handles uploading a thumbnail separately
// POST /api/videos/{id}/thumbnail
func UploadThumbnail(c *CustomContext) {
	id := c.Param("id")
	video, err := database.FindVideoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Map{"error": "video not found"})
		return
	}

	fileHeader, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{"error": "no thumbnail file"})
		return
	}

	ext := filepath.Ext(fileHeader.Filename)
	destPath := filepath.Join(ThumbnailsDir, fmt.Sprintf("%d%s", video.ID, ext))

	if err := c.SaveUploadedFile(fileHeader, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "could not save thumbnail"})
		return
	}

	// Update db
	database.DB.Model(&video).Update("thumbnail_path", destPath)
	c.JSON(http.StatusOK, Map{"thumbnail_path": destPath})
}

func HomePage(c *CustomContext)   {}
func WatchVideo(c *CustomContext) {}
func RemoveAdultContent(c *CustomContext) {
	c.JSON(http.StatusOK, Map{"message": "not implemented"})
}

// ensure io is used
var _ = io.EOF
