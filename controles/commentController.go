package controles

import (
	"net/http"
	"strconv"
	"time"

	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/models"
)

// GetVideoCommentsList returns comments for a video (JSON)
// GET /api/videos/{id}/comments
func GetVideoCommentsList(c *CustomContext) {
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

// CreateComment adds a comment to a video (JSON)
// POST /api/videos/{id}/comments
// Body JSON: { "content": "..." }
func CreateComment(c *CustomContext) {
	id := c.Param("id")
	videoID64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{"error": "invalid video id"})
		return
	}

	var body struct {
		Content string `json:"content"`
		UserID  uint   `json:"user_id"`
	}
	if err := c.Bind(&body); err != nil || body.Content == "" {
		c.JSON(http.StatusBadRequest, Map{"error": "content is required"})
		return
	}

	userID := body.UserID
	if userID == 0 {
		userID = 1 // fallback until auth context is wired
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
