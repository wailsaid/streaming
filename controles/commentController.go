package controles

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/models"
)

// AddComment handles the creation of a new comment via HTMX
func AddComment(c *CustomContext) {
	videoIDStr := c.Query("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 32)
	if err != nil {
		http.Error(c.Writer, "Invalid video ID", http.StatusBadRequest)
		return
	}

	commentContent := c.PostForm("comment")
	if commentContent == "" {
		http.Error(c.Writer, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	// In a real app, get the user ID from the session
	// For now we'll use a placeholder user ID
	userID := uint(1)

	comment := &models.Comment{
		Content:   commentContent,
		VideoID:   uint(videoID),
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	err = database.CreateComment(comment)
	if err != nil {
		http.Error(c.Writer, "Failed to save comment", http.StatusInternalServerError)
		return
	}

	// Get the user to populate the comment response
	user, err := database.FindUserByID(userID)
	if err == nil {
		comment.User = user
	}

	// Render a single comment HTML
	renderSingleComment(c, comment)
}

// GetVideoComments returns comments for a video
func GetVideoComments(c *CustomContext) {
	videoIDStr := c.Query("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 32)
	if err != nil {
		http.Error(c.Writer, "Invalid video ID", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10 // Default limit
	}

	comments := database.GetCommentsByVideoID(uint(videoID), limit)

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(http.StatusOK)

	if len(comments) == 0 {
		c.Writer.Write([]byte("<div class='text-center text-gray-500 py-8'>Be the first to comment!</div>"))
		return
	}

	// Create and parse the template for a single comment
	tmpl := `
		<div class="flex gap-3 pb-4 mb-4 border-b border-gray-100">
			<div class="flex-shrink-0">
				<div class="w-10 h-10 rounded-full bg-gray-200 overflow-hidden">
					<img src="./assets/placeholder-user.jpeg" alt="Profile" class="w-full h-full object-cover">
				</div>
			</div>
			<div class="flex-1">
				<div class="flex items-center gap-2 mb-1">
					<h4 class="font-medium text-gray-900">{{if .User.Username}}{{.User.Username}}{{else}}Anonymous{{end}}</h4>
					<span class="text-xs text-gray-500">{{formatTime .CreatedAt}}</span>
				</div>
				<p class="text-gray-800">{{.Content}}</p>
				<div class="flex gap-4 mt-2 text-sm">
					<button class="text-gray-500 hover:text-blue-600 transition-colors">
						<span class="flex items-center gap-1">
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
								<path d="M7 10v12"/>
								<path d="M15 5.88 14 10h5.83a2 2 0 0 1 1.92 2.56l-2.33 8A2 2 0 0 1 17.5 22H4a2 2 0 0 1-2-2v-8a2 2 0 0 1 2-2h2.76a2 2 0 0 0 1.79-1.11L12 2h0a3.13 3.13 0 0 1 3 3.88Z"/>
							</svg>
							Like
						</span>
					</button>
					<button class="text-gray-500 hover:text-blue-600 transition-colors">Reply</button>
				</div>
			</div>
		</div>
	`

	// Create a template with a function to format the date
	t, err := template.New("comment").Funcs(template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
	}).Parse(tmpl)

	if err != nil {
		http.Error(c.Writer, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render each comment
	for _, comment := range comments {
		err = t.Execute(c.Writer, comment)
		if err != nil {
			http.Error(c.Writer, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Helper function to render a single comment
func renderSingleComment(c *CustomContext, comment *models.Comment) {
	tmpl := `
		<div class="flex gap-3 pb-4 mb-4 border-b border-gray-100">
			<div class="flex-shrink-0">
				<div class="w-10 h-10 rounded-full bg-gray-200 overflow-hidden">
					<img src="./assets/placeholder-user.jpeg" alt="Profile" class="w-full h-full object-cover">
				</div>
			</div>
			<div class="flex-1">
				<div class="flex items-center gap-2 mb-1">
					<h4 class="font-medium text-gray-900">{{if .User.Username}}{{.User.Username}}{{else}}Anonymous{{end}}</h4>
					<span class="text-xs text-gray-500">Just now</span>
				</div>
				<p class="text-gray-800">{{.Content}}</p>
				<div class="flex gap-4 mt-2 text-sm">
					<button class="text-gray-500 hover:text-blue-600 transition-colors">
						<span class="flex items-center gap-1">
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
								<path d="M7 10v12"/>
								<path d="M15 5.88 14 10h5.83a2 2 0 0 1 1.92 2.56l-2.33 8A2 2 0 0 1 17.5 22H4a2 2 0 0 1-2-2v-8a2 2 0 0 1 2-2h2.76a2 2 0 0 0 1.79-1.11L12 2h0a3.13 3.13 0 0 1 3 3.88Z"/>
							</svg>
							Like
						</span>
					</button>
					<button class="text-gray-500 hover:text-blue-600 transition-colors">Reply</button>
				</div>
			</div>
		</div>
	`

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(http.StatusOK)

	t, err := template.New("singleComment").Parse(tmpl)
	if err != nil {
		fmt.Fprintf(c.Writer, "Error: %v", err)
		return
	}

	// Clear form after successful submission
	c.Writer.Header().Set("HX-Trigger", `{"resetForm": ""}`)

	err = t.Execute(c.Writer, comment)
	if err != nil {
		fmt.Fprintf(c.Writer, "Error: %v", err)
		return
	}
}
