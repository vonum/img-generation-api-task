package api

import (
	"net/http"

	"github.com/tutti-ch/backend-coding-task-template/image"
)

type ImageServer struct {
  basePath string
  c chan <- image.Job
}

func NewImageServer(basePath string, c chan <- image.Job) *ImageServer {
  return &ImageServer{basePath, c}
}

func (is *ImageServer) Run(port string) error {
  imageHandler := image.NewImageHandler(image.MaxImageSize, is.c)

  http.HandleFunc("/health", health)
  http.HandleFunc("/upload", imageHandler.Rescale)

  err := http.ListenAndServe(port, nil)
  return err
}

func health(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}
