package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/saidwail/streaming/controles"
	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/env"
	"github.com/saidwail/streaming/midelware"
)

func main() {
	// Initialize app config (loads .env + DB init)
	env.Init()
	database.Connect()

	// Ensure storage directories exist
	for _, dir := range []string{"./tmp", "./storage/videos", "./storage/thumbnails"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("failed to create directory %s: %v", dir, err)
		}
	}

	server := http.NewServeMux()

	// ─── Static Assets ────────────────────────────────────────────
	// Serve Vue dist (production build)
	fileServer := http.FileServer(http.Dir("./webui/dist/assets"))
	server.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	// ─── Auth Routes (public) ────────────────────────────────────
	server.HandleFunc("POST /api/auth/register", adaptHandler(controles.SignUp))
	server.HandleFunc("POST /api/auth/login", adaptHandler(controles.Login))

	// ─── Protected API Routes ────────────────────────────────────
	// Videos
	server.Handle("GET /api/videos", midelware.JwtFilter(adaptHandler(controles.ListVideos)))
	server.Handle("GET /api/videos/search", midelware.JwtFilter(adaptHandler(controles.SearchVideos)))
	server.Handle("GET /api/videos/more", midelware.JwtFilter(adaptHandler(controles.LoadMoreVideos)))
	server.Handle("POST /api/videos/upload", midelware.JwtFilter(adaptHandler(controles.UploadChunk)))

	// Individual video — simulate path param via query for stdlib pattern
	server.Handle("GET /api/videos/{id}", midelware.JwtFilter(adaptHandler(func(c *controles.CustomContext) {
		c.Params["id"] = c.Request.PathValue("id")
		controles.GetVideo(c)
	})))
	server.Handle("GET /api/videos/{id}/recommended", midelware.JwtFilter(adaptHandler(func(c *controles.CustomContext) {
		c.Params["id"] = c.Request.PathValue("id")
		controles.GetRecommendedVideos(c)
	})))
	server.Handle("GET /api/videos/{id}/comments", midelware.JwtFilter(adaptHandler(func(c *controles.CustomContext) {
		c.Params["id"] = c.Request.PathValue("id")
		controles.GetVideoCommentsList(c)
	})))
	server.Handle("POST /api/videos/{id}/comments", midelware.JwtFilter(adaptHandler(func(c *controles.CustomContext) {
		c.Params["id"] = c.Request.PathValue("id")
		controles.CreateComment(c)
	})))
	server.Handle("POST /api/videos/{id}/thumbnail", midelware.JwtFilter(adaptHandler(func(c *controles.CustomContext) {
		c.Params["id"] = c.Request.PathValue("id")
		controles.UploadThumbnail(c)
	})))

	// Streaming — also protected
	server.Handle("GET /api/stream/{id}", midelware.JwtFilter(adaptHandler(func(c *controles.CustomContext) {
		c.Params["id"] = c.Request.PathValue("id")
		controles.StreamVideo(c)
	})))

	// Thumbnail — public (easier for <img> tags)
	server.HandleFunc("GET /api/thumbnail/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := controles.NewCustomContext(w, r)
		ctx.Params["id"] = r.PathValue("id")
		controles.ServeThumbnail(ctx)
	})

	// ─── SPA Fallback ─────────────────────────────────────────────
	// All non-API routes serve the Vue SPA
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Don't intercept API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./webui/dist/index.html")
	})

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(server)))
}

// corsMiddleware adds CORS headers (permissive for development)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// adaptHandler converts a controller function to a standard http.HandlerFunc
func adaptHandler(controller func(*controles.CustomContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse body for POST/PUT
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			ct := r.Header.Get("Content-Type")
			if strings.Contains(ct, "multipart/form-data") {
				r.ParseMultipartForm(100 << 20) // 100MB
			} else if strings.Contains(ct, "application/x-www-form-urlencoded") {
				r.ParseForm()
			}
		}
		ctx := controles.NewCustomContext(w, r)
		controller(ctx)
	}
}

// jsonError writes a JSON error response
func jsonError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
