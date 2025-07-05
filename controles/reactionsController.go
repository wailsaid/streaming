package controles

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// LikeVideo handles a like request via HTMX
func LikeVideo(c *CustomContext) {
	videoIDStr := c.Query("video_id")
	_, err := strconv.ParseUint(videoIDStr, 10, 32)
	if err != nil {
		http.Error(c.Writer, "Invalid video ID", http.StatusBadRequest)
		return
	}

	// In a real application, you would:
	// 1. Check if the user is authenticated
	// 2. Update the like count in the database
	// 3. Return the updated count

	// For demo purposes, we'll just return a random number
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	likeCount := 12000 + r.Intn(2000)

	// Format the number with K suffix
	formattedCount := fmt.Sprintf("%.1fK", float64(likeCount)/1000)

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write([]byte(formattedCount))
}

// DislikeVideo handles a dislike request via HTMX
func DislikeVideo(c *CustomContext) {
	videoIDStr := c.Query("video_id")
	_, err := strconv.ParseUint(videoIDStr, 10, 32)
	if err != nil {
		http.Error(c.Writer, "Invalid video ID", http.StatusBadRequest)
		return
	}

	// In a real application, you would:
	// 1. Check if the user is authenticated
	// 2. Update the dislike count in the database
	// 3. Return the updated count

	// For demo purposes, we'll just return a random number
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	dislikeCount := 100 + r.Intn(200)

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write([]byte(fmt.Sprintf("%d", dislikeCount)))
}
