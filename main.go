package main

import (
	"log"
	"net/http"

	"github.com/saidwail/streaming/controles"
	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/env"
	"github.com/saidwail/streaming/utils"
)

func main() {
	env.Init()
	database.Init()
	database.Connect()

	// Initialize templates
	//templates := template.Must(template.ParseGlob("templ/*.html"))

	// Create a new HTTP server mux
	mux := http.NewServeMux()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./templates/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	// Set up HTTP handlers (converting from Gin to standard library)
	mux.HandleFunc("/", adaptHandler(controles.HomePage))

	mux.HandleFunc("/watch", adaptHandler(controles.WatchVideo))
	mux.HandleFunc("/stream", adaptHandler(controles.StreamVideo))

	mux.HandleFunc("/login", methodRouter(
		adaptHandler(controles.LoginPage),
		adaptPostHandler(controles.Login)))

	mux.HandleFunc("/signup", methodRouter(
		adaptHandler(controles.SignupPage),
		adaptPostHandler(controles.SignUp)))

	mux.HandleFunc("/upload", methodRouter(
		adaptHandler(controles.UploadPage),
		adaptPostHandler(controles.UploadVideo)))

	mux.HandleFunc("/video-list", adaptHandler(controles.ListVideos))

	mux.HandleFunc("/thumbnail", adaptHandler(controles.ServeThumbnail))

	// HTMX endpoints
	//mux.HandleFunc("/api/load-more-videos", adaptHandler(controles.LoadMoreVideos))
	mux.HandleFunc("/api/search-videos", adaptHandler(controles.SearchVideos))
	mux.HandleFunc("/api/comments", methodRouter(
		adaptHandler(controles.GetVideoComments),
		adaptPostHandler(controles.AddComment)))
	//mux.HandleFunc("/api/video-preview", adaptHandler(controles.VideoPreview))
	//mux.HandleFunc("/api/like-video", adaptPostHandler(controles.LikeVideo))
	//mux.HandleFunc("/api/dislike-video", adaptPostHandler(controles.DislikeVideo))

	// Initialize MinIO client
	if err := utils.InitMinioClient(); err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Start the server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// methodRouter routes requests based on HTTP method
func methodRouter(getHandler, postHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getHandler(w, r)
		case http.MethodPost:
			postHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// adaptHandler converts controllers to standard http handlers
func adaptHandler(controller func(*controles.CustomContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom context
		ctx := controles.NewCustomContext(w, r)
		controller(ctx)
	}
}

// adaptPostHandler for POST requests
func adaptPostHandler(controller func(*controles.CustomContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data first
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}
		}

		// Create a custom context
		ctx := controles.NewCustomContext(w, r)
		controller(ctx)
	}
}
