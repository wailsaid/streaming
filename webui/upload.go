package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./dist")))

	http.HandleFunc("POST /upload", uploadHander)
	log.Println("starting server...")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func uploadHander(w http.ResponseWriter, r *http.Request) {

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
	video, err := os.Create("tmp/" + filename + ".mp4")
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

}
