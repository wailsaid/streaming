package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/saidwail/streaming/controles"
)

func main() {
	/* env.Init()
	database.Init()
	database.Connect() */

	// Create a new HTTP server mux
	server := http.NewServeMux()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./webui/dist/assets"))
	server.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	//Default route
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./webui/dist/index.html")
	})

	server.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		dir, err := os.ReadDir("tmp/videos")
		if err != nil {
			log.Fatal(err)
		}
		videos := make([]map[string]interface{}, 0)
		for _, file := range dir {
			video := map[string]interface{}{}
			v, _ := file.Info()
			video["title"] = v.Name()
			video["path"] = "tmp/videos/" + v.Name()
			videos = append(videos, video)
		}

		w.Header().Add("Content-type", "application/json")

		if err := json.NewEncoder(w).Encode(videos); err != nil {

		}

		//w.Write([]byte(`{"message": "Search endpoint is not implemented yet"}`))

	})

	server.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		file, fileHeader, _ := r.FormFile("file")

		filename := r.FormValue("file_name")
		index, _ := strconv.Atoi(r.FormValue("index"))
		total, _ := strconv.Atoi(r.FormValue("total"))

		log.Printf("reciving file %v part %v size %v total %v", filename, index, fileHeader.Size, total)
		//filename := r.FormValue("file")

		fw, _ := os.Create("tmp/" + filename + "_" + strconv.Itoa(index) + ".tmp")
		b := make([]byte, fileHeader.Size)
		file.Read(b)
		fw.Write(b)
		defer fw.Close()

		if index+1 != total {
			return
		}
		log.Printf("all parts received for %v, merging...", filename)
		w.WriteHeader(201)
		w.Write([]byte(`{"message": "File uploaded successfully"}`))

		go func() {
			video, err := os.Create("tmp/videos/" + filename)
			if err != nil {
				log.Fatal(err)
			}
			defer video.Close()

			for i := range total {
				b, err := os.ReadFile("tmp/" + filename + "_" + strconv.Itoa(i) + ".tmp")
				if err != nil {
					log.Fatal(err)
				}
				video.Write(b)
				err = os.Remove("tmp/" + filename + "_" + strconv.Itoa(i) + ".tmp")
				if err != nil {
					log.Fatal(err)

				}
			}
			log.Printf("File %v merged successfully", filename)
		}()

	})

	/*
		server.HandleFunc("/watch", adaptHandler(controles.WatchVideo))
			server.HandleFunc("/stream", adaptHandler(controles.StreamVideo))

			server.HandleFunc("/login", methodRouter(
				adaptHandler(controles.LoginPage),
				adaptPostHandler(controles.Login)))

			server.HandleFunc("/signup", methodRouter(
				adaptHandler(controles.SignupPage),
				adaptPostHandler(controles.SignUp)))

			server.HandleFunc("/upload", methodRouter(
				adaptHandler(controles.UploadPage),
				adaptPostHandler(controles.UploadVideo)))

			server.HandleFunc("/video-list", adaptHandler(controles.ListVideos))

			server.HandleFunc("/thumbnail", adaptHandler(controles.ServeThumbnail))

			// HTMX endpoints
			//mux.HandleFunc("/api/load-more-videos", adaptHandler(controles.LoadMoreVideos))
			server.HandleFunc("/api/search-videos", adaptHandler(controles.SearchVideos))
			server.HandleFunc("/api/comments", methodRouter(
				adaptHandler(controles.GetVideoComments),
				adaptPostHandler(controles.AddComment))) */
	//mux.HandleFunc("/api/video-preview", adaptHandler(controles.VideoPreview))
	//mux.HandleFunc("/api/like-video", adaptPostHandler(controles.LikeVideo))
	//mux.HandleFunc("/api/dislike-video", adaptPostHandler(controles.DislikeVideo))

	// Initialize MinIO client
	/* if err := utils.InitMinioClient(); err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	} */

	// Start the server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}

func adaptHandler(controller func(*controles.CustomContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := controles.NewCustomContext(w, r)
		controller(ctx)
	}
}

/* // methodRouter routes requests based on HTTP method
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
*/
