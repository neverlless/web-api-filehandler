package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	port          string
	filesDir      string
	ipAddress     string
	maxUploadSize int64
)

func main() {
	flag.StringVar(&port, "port", "8080", "Port for launching the web server")
	flag.StringVar(&filesDir, "dir", "./files", "Directory for saving files")
	flag.StringVar(&ipAddress, "ip", "0.0.0.0", "IP address for launching the server")
	flag.Int64Var(&maxUploadSize, "max-upload-size", 10<<20, "Maximum upload file size in bytes")
	flag.Parse()

	// Logger setup
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}
	log.Out = os.Stdout

	// Checking and creating directory for files
	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(filesDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	log.Println("Starting server...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequests(w, r, log)
	})

	fullAddress := fmt.Sprintf("%s:%s", ipAddress, port)
	log.Printf("Server running on http://%s", fullAddress)
	log.Printf("Upload endpoint: http://%s/api/filehandler/", fullAddress)

	if err := http.ListenAndServe(fullAddress, nil); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}

func handleRequests(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	log.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"ip":     r.RemoteAddr,
	}).Info("Request received")

	// Modifying the condition to account for all requests starting with /api/filehandler/
	if r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/api/filehandler/") {
		uploadFileHandler(w, r, log)
		return
	}

	filePath := filepath.Join(filesDir, filepath.Clean(r.URL.Path))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Limiting the request body size to prevent memory overflow
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.WithFields(logrus.Fields{
			"ip":  r.RemoteAddr,
			"err": err,
		}).Error("Maximum file size exceeded")
		http.Error(w, "Maximum file size exceeded", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.WithFields(logrus.Fields{
			"ip":  r.RemoteAddr,
			"err": err,
		}).Error("Error receiving file")
		http.Error(w, "Error receiving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Extracting subpath from URL
	subPath := strings.TrimPrefix(r.URL.Path, "/api/filehandler/")
	if subPath == "" || strings.HasSuffix(subPath, "/") {
		subPath = filepath.Join(subPath, handler.Filename) // Adding file name if subpath is empty or ends with /
	}

	filePath := filepath.Join(filesDir, filepath.Clean(subPath))
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm) // Creating all necessary directories

	dst, err := os.Create(filePath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"path": filePath,
			"err":  err,
		}).Error("Error creating file")
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		log.WithFields(logrus.Fields{
			"path": filePath,
			"err":  err,
		}).Error("Error saving file")
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	log.WithFields(logrus.Fields{
		"path": filePath,
		"ip":   r.RemoteAddr,
	}).Info("File successfully uploaded")

	fmt.Fprintf(w, "File %s uploaded successfully", handler.Filename)
}
