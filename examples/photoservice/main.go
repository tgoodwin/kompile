package main

import (
	"fmt"
	"image"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

const standardWidth = 100
const standardHeight = 100

func resizeImage(file multipart.File, fileName string, errCh chan error) {
	fmt.Println("Processing image:", fileName)
	img, _, err := image.Decode(file)
	if err != nil {
		errCh <- fmt.Errorf("decoding the image: %v", err)
	}

	resizedImage := imaging.Resize(img, standardWidth, standardHeight, imaging.Lanczos)

	// write out the processed image
	outpath := "processed_" + fileName
	outfile, err := os.Create("./" + outpath)
	if err != nil {
		errCh <- fmt.Errorf("creating the processed image: %v", err)
	}
	fmt.Printf("Saving processed image to: %s\n", outpath)

	defer outfile.Close()

	// encode the resized image and save it
	err = imaging.Encode(outfile, resizedImage, imaging.JPEG)
	if err != nil {
		fmt.Printf("saving the processed image: %v", err)
		errCh <- fmt.Errorf("saving the processed image: %v", err)
	}

	errCh <- nil
}

type Server struct {
	jobStatus map[string]string
	mu        sync.Mutex
}

func NewServer() *Server {
	return &Server{
		jobStatus: make(map[string]string),
	}
}

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	// limit the maximum upload size to 10MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded file: %+v/n", handler.Filename)

	// process the file in memory for now

	// Generate a new job ID
	jobID := uuid.New().String()

	// Save the initial job status
	s.mu.Lock()
	s.jobStatus[jobID] = "processing"
	s.mu.Unlock()

	// Dummy response for now (actual processing will happen asynchronously)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"job_id": "` + jobID + `"}`))

	// TODO handle updating job status upon completion or error
	errChan := make(chan error)
	go resizeImage(file, handler.Filename, errChan)
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "job_id")

	s.mu.Lock()
	status, exists := s.jobStatus[jobID]
	s.mu.Unlock()

	if !exists {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(`{"job_id": "` + jobID + `", "status": "` + status + `"}`))
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	server := NewServer()
	r.Post("/upload", server.uploadHandler)
	r.Get("/status/{job_id}", server.statusHandler)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
